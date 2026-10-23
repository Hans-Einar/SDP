package runtime

// ReplaceCollection is an explicit application refresh; root becomes loaded.
func (s *Session) ReplaceCollection(t CollectionTarget, d CollectionData) error {
	if err := s.ValidateCollectionTarget(t); err != nil {
		return err
	}
	n := s.copyState()
	c := n.collections[t.Handle.Path]
	if err := c.replace(t, d, true); err != nil {
		return collectionFault(c.InstancePath, err)
	}
	return s.publish(n)
}
func (s *Session) ReplaceChildren(t CollectionTarget, d CollectionData) error {
	if err := s.ValidateCollectionTarget(t); err != nil {
		return err
	}
	n := s.copyState()
	c := n.collections[t.Handle.Path]
	if err := c.replace(t, d, false); err != nil {
		return collectionFault(c.InstancePath, err)
	}
	return s.publish(n)
}
func (c *collection) descendant(id, parent ItemID) bool {
	x, ok := c.item(id)
	for ok && x.Parent != "" {
		if x.Parent == parent {
			return true
		}
		x, ok = c.item(x.Parent)
	}
	return false
}
func (c *collection) replace(t CollectionTarget, d CollectionData, whole bool) error {
	if len(d.Items) > 4096 {
		return fault("collection-limit", "Collection exceeds 4096 items")
	}
	old := copyCollection(c)
	data := copyData(d)
	if !whole && t.ItemID != "" {
		parent, ok := c.item(t.ItemID)
		if !ok || c.Handle.Kind != "tree" || !parent.HasChildren {
			return fault("collection-item", "Child replacement requires a tree branch")
		}
		added := map[ItemID]CollectionItem{}
		for _, x := range d.Items {
			if _, ok := added[x.ID]; ok {
				return fault("collection-data", "Duplicate replacement ID")
			}
			added[x.ID] = x
		}
		for _, x := range d.Items {
			seen := map[ItemID]bool{x.ID: true}
			p := x.Parent
			for p != t.ItemID {
				if p == "" || seen[p] {
					return fault("collection-data", "Replacement item is not within the target subtree")
				}
				seen[p] = true
				ancestor, ok := added[p]
				if !ok {
					return fault("collection-data", "Replacement parent is outside supplied subtree")
				}
				p = ancestor.Parent
			}
		}
		data.Items = nil
		for _, x := range c.Data.Items {
			if c.descendant(x.ID, t.ItemID) {
				continue
			}
			if x.ID == t.ItemID {
				x.ChildrenLoaded = true
			}
			data.Items = append(data.Items, x)
		}
		data.Items = append(data.Items, d.Items...)
	}
	if err := validateData(c.Handle.Kind, data, true, c.provider.Load != nil); err != nil {
		return err
	}
	c.Data = data
	c.RootLoaded = true
	c.AutoLoadPending = false
	c.Generation++
	c.Request = nil
	c.resetStatus(old)
	for id := range c.Expanded {
		x, ok := c.item(id)
		previous, _ := old.item(id)
		if !ok || !x.HasChildren || x.Kind != previous.Kind {
			delete(c.Expanded, id)
		}
	}
	if x, ok := c.item(c.Selected); !ok || x.Kind != Row {
		c.Selected = ""
	}
	if x, ok := c.item(c.Focused); !ok || x.Kind == Separator {
		c.Focused = ""
	} else if previous, _ := old.item(c.Focused); previous.Kind != x.Kind {
		c.Focused = ""
	}
	return nil
}

func (s *Session) interactive(t CollectionTarget) (*collection, CollectionItem, error) {
	if err := s.ValidateCollectionTarget(t); err != nil {
		return nil, CollectionItem{}, err
	}
	w := s.widgets[t.Handle.Path]
	c := s.collections[t.Handle.Path]
	if !w.Enabled || !w.Visible {
		return nil, CollectionItem{}, fault("inactive-widget", w.Handle.Path)
	}
	x, _ := c.item(t.ItemID)
	if t.ItemID != "" && !c.visible(t.ItemID) {
		return nil, CollectionItem{}, fault("hidden-item", "Collection item is collapsed")
	}
	return c, x, nil
}
func (s *Session) SelectItem(t CollectionTarget) error {
	_, x, err := s.interactive(t)
	if err != nil {
		return err
	}
	if t.ItemID == "" || x.Kind != Row {
		return fault("collection-item", "Selection requires a row")
	}
	n := s.copyState()
	n.collections[t.Handle.Path].Selected = t.ItemID
	return s.publish(n)
}
func (s *Session) FocusItem(t CollectionTarget) error {
	if !s.inputAllowed(s.widgets[t.Handle.Path]) {
		return fault("focus", "Collection is not focusable")
	}
	_, x, err := s.interactive(t)
	if err != nil {
		return err
	}
	if t.ItemID != "" && x.Kind == Separator {
		return fault("collection-item", "Separator is not focusable")
	}
	n := s.copyState()
	n.collections[t.Handle.Path].Focused = t.ItemID
	n.focused = t.Handle.Path
	n.rememberFocus(t.Handle.Path)
	return s.publish(n)
}
func (s *Session) ExpandItem(t CollectionTarget) error {
	c, x, err := s.interactive(t)
	if err != nil {
		return err
	}
	if t.ItemID == "" || !x.HasChildren || x.Kind == Separator {
		return fault("collection-item", "Expansion requires a branch")
	}
	n := s.copyState()
	n.collections[t.Handle.Path].Expanded[t.ItemID] = true
	if !x.ChildrenLoaded && c.Status[t.ItemID].Phase != LoadError && c.Status[t.ItemID].Phase != Loading {
		if _, err = n.begin(t); err != nil {
			return err
		}
	}
	return s.publish(n)
}
func (s *Session) CollapseItem(t CollectionTarget) error {
	_, x, err := s.interactive(t)
	if err != nil {
		return err
	}
	if t.ItemID == "" || !x.HasChildren {
		return fault("collection-item", "Collapse requires a branch")
	}
	n := s.copyState()
	c := n.collections[t.Handle.Path]
	delete(c.Expanded, t.ItemID)
	if c.Request != nil && (c.Request.Target.ItemID == t.ItemID || c.descendant(c.Request.Target.ItemID, t.ItemID)) {
		c.cancel()
	}
	if c.descendant(c.Focused, t.ItemID) {
		c.Focused = t.ItemID
	}
	return s.publish(n)
}
func (s *Session) RetryItem(t CollectionTarget) error {
	c, x, err := s.interactive(t)
	if err != nil {
		return err
	}
	if t.ItemID != "" && (!x.HasChildren || x.Kind == Separator) {
		return fault("collection-item", "Retry requires a branch or root")
	}
	if c.Request != nil && c.Request.Target == t {
		return nil
	}
	status := c.Status[t.ItemID]
	unloaded := !x.ChildrenLoaded
	if t.ItemID == "" {
		unloaded = !c.RootLoaded
	}
	// A loading root becomes unloaded on compatible reload, while its one-shot
	// auto flag stays consumed. Explicit Load recovery must remain possible.
	pausedRoot := t.ItemID == "" && unloaded && !c.AutoLoadPending && status.Phase == Unloaded
	if status.Phase != LoadError && !(status.Phase == Canceled && unloaded) && !pausedRoot {
		return fault("collection-load", "Target is not recoverable")
	}
	n := s.copyState()
	if t.ItemID != "" {
		n.collections[t.Handle.Path].Expanded[t.ItemID] = true
	}
	if _, err = n.begin(t); err != nil {
		return err
	}
	return s.publish(n)
}
