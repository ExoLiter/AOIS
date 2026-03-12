package hashtable

func NewTable(size int) (*Table, error) {
	if size < MinTableSize {
		return nil, ErrTableSize
	}
	slots := make([]Entry, size)
	for i := range slots {
		slots[i] = emptyEntry()
	}
	return &Table{
		size:  size,
		base:  BaseAddress,
		slots: slots,
	}, nil
}

func (t *Table) Size() int {
	return t.size
}

func (t *Table) Insert(key string, value string) error {
	normalized, v, h, err := t.prepareKey(key)
	if err != nil {
		return err
	}
	if _, ok := t.findIndex(normalized, h); ok {
		return ErrDuplicateKey
	}
	idx, ok := t.findFreeSlot(h)
	if !ok {
		return ErrTableFull
	}
	t.writeEntry(idx, normalized, value, v, h)
	t.rebuildChains()
	return nil
}

func (t *Table) Find(key string) (Entry, bool) {
	normalized, v, h, err := t.prepareKey(key)
	if err != nil {
		return Entry{}, false
	}
	idx, ok := t.findIndex(normalized, h)
	if !ok {
		return Entry{}, false
	}
	entry := t.slots[idx]
	entry.V = v
	entry.Home = h
	return entry, true
}

func (t *Table) Update(key string, value string) error {
	normalized, _, h, err := t.prepareKey(key)
	if err != nil {
		return err
	}
	idx, ok := t.findIndex(normalized, h)
	if !ok {
		return ErrNotFound
	}
	t.slots[idx].Value = value
	return nil
}

func (t *Table) Delete(key string) error {
	normalized, _, h, err := t.prepareKey(key)
	if err != nil {
		return err
	}
	idx, ok := t.findIndex(normalized, h)
	if !ok {
		return ErrNotFound
	}
	t.markDeleted(idx)
	t.rebuildChains()
	return nil
}

func (t *Table) LoadFactor() float64 {
	occupied := 0
	for i := range t.slots {
		if t.slots[i].Flags.Occupied {
			occupied++
		}
	}
	return float64(occupied) / float64(t.size)
}

func (t *Table) Entries() []Entry {
	entries := make([]Entry, 0, t.size)
	for i := range t.slots {
		entries = append(entries, t.slots[i])
	}
	return entries
}

func (t *Table) prepareKey(key string) (string, int, int, error) {
	normalized, err := normalizeKey(key)
	if err != nil {
		return "", 0, 0, err
	}
	v, err := computeV(normalized)
	if err != nil {
		return "", 0, 0, err
	}
	h := computeH(v, t.size, t.base)
	return normalized, v, h, nil
}

func (t *Table) findFreeSlot(start int) (int, bool) {
	for offset := 0; offset < t.size; offset++ {
		idx := (start + offset) % t.size
		if !t.slots[idx].Flags.Occupied {
			return idx, true
		}
	}
	return 0, false
}

func (t *Table) findIndex(normalized string, start int) (int, bool) {
	for offset := 0; offset < t.size; offset++ {
		idx := (start + offset) % t.size
		slot := t.slots[idx]
		if slot.Flags.Occupied && slot.Key == normalized {
			return idx, true
		}
		if !slot.Flags.Occupied && !slot.Flags.Deleted {
			return 0, false
		}
	}
	return 0, false
}

func (t *Table) writeEntry(idx int, key string, value string, v int, h int) {
	t.slots[idx] = Entry{
		Key:   key,
		Value: value,
		V:     v,
		Home:  h,
		Flags: Flags{
			Occupied: true,
		},
		Next: NoNextIndex,
	}
}

func (t *Table) markDeleted(idx int) {
	slot := &t.slots[idx]
	slot.Flags.Occupied = false
	slot.Flags.Deleted = true
	slot.Flags.Collision = false
	slot.Flags.Terminal = false
	slot.Next = NoNextIndex
}

func emptyEntry() Entry {
	return Entry{Next: NoNextIndex}
}
