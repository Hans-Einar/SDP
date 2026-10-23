package runtime

func (c *collection) cancel() {
	if c.Request == nil {
		return
	}
	id := c.Request.Target.ItemID
	c.Request = nil
	c.Status[id] = LoadStatus{Phase: Canceled}
}

// begin mutates only a detached candidate and never invokes a provider.
func (s *Session) begin(t CollectionTarget) (LoadRequest, error) {
	c, x, err := s.interactive(t)
	if err != nil {
		return LoadRequest{}, err
	}
	if c.provider.Load == nil {
		return LoadRequest{}, fault("collection-load", "No loader bound")
	}
	if t.ItemID != "" && (!x.HasChildren || !c.Expanded[t.ItemID]) {
		return LoadRequest{}, fault("collection-load", "Load requires an expanded branch")
	}
	if c.Request != nil && c.Request.Target == t {
		return *c.Request, nil
	}
	c.cancel()
	c.requestID++
	r := LoadRequest{Target: t, RequestID: c.requestID, ProviderEpoch: c.ProviderEpoch}
	c.Request = &r
	c.Status[t.ItemID] = LoadStatus{Phase: Loading}
	if t.ItemID == "" {
		c.AutoLoadPending = false
	}
	return r, nil
}
func (s *Session) BeginLoad(t CollectionTarget) (LoadRequest, error) {
	if err := s.ValidateCollectionTarget(t); err != nil {
		return LoadRequest{}, err
	}
	n := s.copyState()
	r, err := n.begin(t)
	if err != nil {
		return LoadRequest{}, err
	}
	if err = s.publish(n); err != nil {
		return LoadRequest{}, err
	}
	return r, nil
}
func (s *Session) CancelLoad(h Handle) error {
	if _, err := s.lookup(h); err != nil {
		return err
	}
	if s.collections[h.Path] == nil {
		return fault("collection", "Collection is not bound")
	}
	if s.collections[h.Path].Request == nil {
		return nil
	}
	n := s.copyState()
	n.collections[h.Path].cancel()
	return s.publish(n)
}
func (s *Session) validRequest(r LoadRequest) error {
	c, _, err := s.interactive(r.Target)
	if err != nil {
		return err
	}
	if c.Request == nil || *c.Request != r || c.ProviderEpoch != r.ProviderEpoch {
		return fault("stale-load", "Load request is no longer current")
	}
	if r.Target.ItemID != "" && !c.Expanded[r.Target.ItemID] {
		return fault("stale-load", "Branch is no longer expanded")
	}
	return nil
}

// CompleteLoad publishes data only after full validation; matching failures become
// recoverable error state. A nil return means the success/error state was accepted.
func (s *Session) CompleteLoad(r LoadRequest, d CollectionData, loadErr error) error {
	if err := s.validRequest(r); err != nil {
		return err
	}
	if loadErr == nil {
		n := s.copyState()
		c := n.collections[r.Target.Handle.Path]
		loadErr = c.replace(r.Target, d, false)
		if loadErr == nil {
			loadErr = s.publish(n)
			if loadErr == nil {
				return nil
			}
		}
	}
	n := s.copyState()
	c := n.collections[r.Target.Handle.Path]
	c.Request = nil
	c.Status[r.Target.ItemID] = LoadStatus{Phase: LoadError, Error: boundedError(loadErr)}
	return s.publish(n)
}
