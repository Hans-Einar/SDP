package runtime

import "context"

type ItemID string
type ItemKind string

const (
	Row       ItemKind = "row"
	Group     ItemKind = "group"
	Separator ItemKind = "separator"
)

type CollectionItem struct {
	ID, Parent                  ItemID
	Kind                        ItemKind
	Label                       string
	HasChildren, ChildrenLoaded bool
}
type CollectionData struct{ Items []CollectionItem }
type CollectionTarget struct {
	Handle                              Handle
	ModelRevision, CollectionGeneration uint64
	ItemID                              ItemID
}
type LoadRequest struct {
	Target                   CollectionTarget
	RequestID, ProviderEpoch uint64
}
type CollectionProvider struct {
	ID         string
	Epoch      uint64
	Initial    CollectionData
	RootLoaded bool
	Load       func(context.Context, LoadRequest) (CollectionData, error)
}
type LoadPhase string

const (
	Unloaded  LoadPhase = "unloaded"
	Loading   LoadPhase = "loading"
	Loaded    LoadPhase = "loaded"
	LoadError LoadPhase = "error"
	Canceled  LoadPhase = "canceled"
)

type LoadStatus struct {
	Phase LoadPhase
	Error string
}
type CollectionState struct {
	Handle                      Handle
	InstancePath                string
	Generation                  uint64
	ProviderID                  string
	ProviderEpoch               uint64
	Data                        CollectionData
	RootLoaded, AutoLoadPending bool
	Selected, Focused           ItemID
	Expanded                    map[ItemID]bool
	Status                      map[ItemID]LoadStatus
	Request                     *LoadRequest
}
type collection struct {
	CollectionState
	provider  CollectionProvider
	requestID uint64
}

func isCollection(kind string) bool { return kind == "tree" || kind == "list" }
func copyData(d CollectionData) CollectionData {
	return CollectionData{Items: append([]CollectionItem(nil), d.Items...)}
}
func copyCollection(c *collection) *collection {
	n := *c
	n.Data = copyData(c.Data)
	n.Expanded = make(map[ItemID]bool, len(c.Expanded))
	for k, v := range c.Expanded {
		n.Expanded[k] = v
	}
	n.Status = make(map[ItemID]LoadStatus, len(c.Status))
	for k, v := range c.Status {
		n.Status[k] = v
	}
	if c.Request != nil {
		r := *c.Request
		n.Request = &r
	}
	n.provider.Initial = copyData(c.provider.Initial)
	return &n
}
func (c *collection) item(id ItemID) (CollectionItem, bool) {
	for _, item := range c.Data.Items {
		if item.ID == id {
			return item, true
		}
	}
	return CollectionItem{}, false
}
func (c *collection) visible(id ItemID) bool {
	item, ok := c.item(id)
	if !ok {
		return false
	}
	for item.Parent != "" {
		if !c.Expanded[item.Parent] {
			return false
		}
		item, ok = c.item(item.Parent)
		if !ok {
			return false
		}
	}
	return true
}
func (s *Session) Collection(h Handle) (CollectionState, bool) {
	if _, err := s.lookup(h); err != nil {
		return CollectionState{}, false
	}
	c := s.collections[h.Path]
	if c == nil {
		return CollectionState{}, false
	}
	return copyCollection(c).CollectionState, true
}
func (s *Session) Provider(h Handle) (CollectionProvider, bool) {
	if _, err := s.lookup(h); err != nil {
		return CollectionProvider{}, false
	}
	c := s.collections[h.Path]
	if c == nil {
		return CollectionProvider{}, false
	}
	p := c.provider
	p.Initial = copyData(p.Initial)
	return p, true
}
func (s *Session) Target(h Handle, id ItemID) (CollectionTarget, error) {
	if _, err := s.lookup(h); err != nil {
		return CollectionTarget{}, err
	}
	c := s.collections[h.Path]
	if c == nil {
		return CollectionTarget{}, fault("collection", "Collection provider is not bound")
	}
	t := CollectionTarget{h, s.Revision, c.Generation, id}
	return t, s.ValidateCollectionTarget(t)
}
func (s *Session) ValidateCollectionTarget(t CollectionTarget) error {
	if _, err := s.lookup(t.Handle); err != nil {
		return err
	}
	c := s.collections[t.Handle.Path]
	if c == nil || !isCollection(t.Handle.Kind) {
		return fault("collection", "Target requires a bound tree or list")
	}
	if t.ModelRevision != s.Revision || t.CollectionGeneration != c.Generation {
		return fault("stale-collection", "Collection target is no longer current")
	}
	if t.ItemID != "" {
		if _, ok := c.item(t.ItemID); !ok {
			return fault("collection-item", "Unknown item identity")
		}
	}
	return nil
}

// BindProviders validates the exact normalized instance-path map without invoking Load.
// Matching ID/epoch retains state; a changed binding resets data and revokes requests.
func (s *Session) BindProviders(providers map[string]CollectionProvider) error {
	if s.closed {
		return fault("closed", "Session is closed")
	}
	n := s.copyState()
	n.collections = map[string]*collection{}
	used := map[string]bool{}
	for _, w := range s.Widgets() {
		if !isCollection(w.Handle.Kind) {
			continue
		}
		p, ok := providers[w.InstancePath]
		if !ok {
			return fault("provider", "Missing provider at "+w.InstancePath)
		}
		used[w.InstancePath] = true
		if !validPlain(p.ID, 1024, false) {
			return fault("provider", "Provider ID required at "+w.InstancePath)
		}
		if err := validateData(w.Handle.Kind, p.Initial, p.RootLoaded, p.Load != nil); err != nil {
			return collectionFault(w.InstancePath, err)
		}
		if old := s.collections[w.Handle.Path]; old != nil && old.ProviderID == p.ID && old.ProviderEpoch == p.Epoch {
			n.collections[w.Handle.Path] = copyCollection(old)
			continue
		}
		c := &collection{CollectionState: CollectionState{Handle: w.Handle, InstancePath: w.InstancePath, Generation: 1, ProviderID: p.ID, ProviderEpoch: p.Epoch, Data: copyData(p.Initial), RootLoaded: p.RootLoaded, AutoLoadPending: !p.RootLoaded, Expanded: map[ItemID]bool{}, Status: map[ItemID]LoadStatus{}}, provider: p}
		c.provider.Initial = copyData(p.Initial)
		if old := s.collections[w.Handle.Path]; old != nil {
			c.Generation = old.Generation + 1
			c.requestID = old.requestID
		}
		c.resetStatus(nil)
		n.collections[w.Handle.Path] = c
	}
	if len(used) != len(providers) {
		return fault("provider", "Unused or wrong-kind provider binding")
	}
	return s.publish(n)
}
