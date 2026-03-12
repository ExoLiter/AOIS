package hashtable

func (t *Table) rebuildChains() {
	resetChainFlags(t.slots)
	groups := groupByHome(t.slots)
	for home, indices := range groups {
		ordered := t.orderByProbe(home, indices)
		t.applyChain(ordered)
	}
}

func resetChainFlags(slots []Entry) {
	for i := range slots {
		slots[i].Flags.Collision = false
		slots[i].Flags.Terminal = false
		slots[i].Next = NoNextIndex
	}
}

func groupByHome(slots []Entry) map[int][]int {
	groups := make(map[int][]int)
	for idx, slot := range slots {
		if slot.Flags.Occupied {
			groups[slot.Home] = append(groups[slot.Home], idx)
		}
	}
	return groups
}

func (t *Table) orderByProbe(home int, indices []int) []int {
	set := indexSet(indices)
	ordered := make([]int, 0, len(indices))
	for offset := 0; offset < t.size; offset++ {
		idx := (home + offset) % t.size
		if set[idx] {
			ordered = append(ordered, idx)
		}
	}
	return ordered
}

func indexSet(indices []int) map[int]bool {
	set := make(map[int]bool)
	for _, idx := range indices {
		set[idx] = true
	}
	return set
}

func (t *Table) applyChain(ordered []int) {
	count := len(ordered)
	for i, idx := range ordered {
		entry := &t.slots[idx]
		entry.Flags.Collision = count > 1
		if count == 1 {
			entry.Flags.Terminal = true
			entry.Next = idx
			continue
		}
		isLast := i == count-1
		entry.Flags.Terminal = isLast
		if isLast {
			entry.Next = idx
			continue
		}
		entry.Next = ordered[i+1]
	}
}
