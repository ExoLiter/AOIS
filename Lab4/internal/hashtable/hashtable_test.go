package hashtable

import (
	"math"
	"testing"
)

func TestComputeVRussian(t *testing.T) {
	v, err := computeV("Вяткин")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	const expected = 98
	if v != expected {
		t.Fatalf("expected %d, got %d", expected, v)
	}
}

func TestComputeVLatin(t *testing.T) {
	v, err := computeV("Test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	const expected = 498
	if v != expected {
		t.Fatalf("expected %d, got %d", expected, v)
	}
}

func TestComputeVInvalid(t *testing.T) {
	_, err := computeV("A")
	if err != ErrKeyInvalid {
		t.Fatalf("expected ErrKeyInvalid, got %v", err)
	}
	_, err = computeV("AБ")
	if err != ErrKeyAlphabet {
		t.Fatalf("expected ErrKeyAlphabet, got %v", err)
	}
}

func TestNewTableSize(t *testing.T) {
	_, err := NewTable(MinTableSize - 1)
	if err != ErrTableSize {
		t.Fatalf("expected ErrTableSize, got %v", err)
	}
}

func TestInsertFindUpdateDelete(t *testing.T) {
	table := newTestTable(t)
	if err := table.Insert("Анатомия", "A"); err != nil {
		t.Fatalf("insert failed: %v", err)
	}
	entry, ok := table.Find("Анатомия")
	if !ok || entry.Value != "A" {
		t.Fatalf("find failed: %+v", entry)
	}
	if err := table.Update("Анатомия", "B"); err != nil {
		t.Fatalf("update failed: %v", err)
	}
	entry, ok = table.Find("Анатомия")
	if !ok || entry.Value != "B" {
		t.Fatalf("update verification failed: %+v", entry)
	}
	if err := table.Delete("Анатомия"); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	_, ok = table.Find("Анатомия")
	if ok {
		t.Fatalf("expected deleted key to be missing")
	}
}

func TestDuplicateKey(t *testing.T) {
	table := newTestTable(t)
	if err := table.Insert("Анатомия", "A"); err != nil {
		t.Fatalf("insert failed: %v", err)
	}
	if err := table.Insert("Анатомия", "B"); err != ErrDuplicateKey {
		t.Fatalf("expected ErrDuplicateKey, got %v", err)
	}
}

func TestUpdateDeleteNotFound(t *testing.T) {
	table := newTestTable(t)
	if err := table.Update("Анатомия", "A"); err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
	if err := table.Delete("Анатомия"); err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestCollisionChain(t *testing.T) {
	table := newTestTable(t)
	words := []string{"Анатомия", "Анализ", "Антисептик"}
	for _, word := range words {
		if err := table.Insert(word, "X"); err != nil {
			t.Fatalf("insert failed: %v", err)
		}
	}
	_, _, h, err := table.prepareKey(words[0])
	if err != nil {
		t.Fatalf("prepareKey failed: %v", err)
	}
	indices := make([]int, 0, len(words))
	for _, word := range words {
		idx, ok := table.findIndex(stringsUpper(word), h)
		if !ok {
			t.Fatalf("expected to find %s", word)
		}
		indices = append(indices, idx)
	}
	ordered := table.orderByProbe(h, indices)
	for i, idx := range ordered {
		slot := table.slots[idx]
		if !slot.Flags.Collision {
			t.Fatalf("expected collision flag on %s", slot.Key)
		}
		isLast := i == len(ordered)-1
		if slot.Flags.Terminal != isLast {
			t.Fatalf("terminal flag mismatch for %s", slot.Key)
		}
		if !isLast && slot.Next != ordered[i+1] {
			t.Fatalf("next pointer mismatch for %s", slot.Key)
		}
		if isLast && slot.Next != idx {
			t.Fatalf("last node should point to itself")
		}
	}
}

func TestDeleteKeepsChainSearch(t *testing.T) {
	table := newTestTable(t)
	if err := table.Insert("Анатомия", "A"); err != nil {
		t.Fatalf("insert failed: %v", err)
	}
	if err := table.Insert("Анализ", "B"); err != nil {
		t.Fatalf("insert failed: %v", err)
	}
	if err := table.Insert("Антисептик", "C"); err != nil {
		t.Fatalf("insert failed: %v", err)
	}
	if err := table.Delete("Анатомия"); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if _, ok := table.Find("Анализ"); !ok {
		t.Fatalf("expected to find Analiz after deletion")
	}
	if _, ok := table.Find("Антисептик"); !ok {
		t.Fatalf("expected to find Antiseptik after deletion")
	}
}

func TestLoadFactor(t *testing.T) {
	table := newTestTable(t)
	for _, key := range []string{"Анатомия", "Анализ", "Антисептик", "Вирус"} {
		if err := table.Insert(key, "X"); err != nil {
			t.Fatalf("insert failed: %v", err)
		}
	}
	lf := table.LoadFactor()
	expected := float64(4) / float64(DefaultTableSize)
	assertFloatNear(t, lf, expected)
}

func TestTableFull(t *testing.T) {
	table := newTestTable(t)
	for i := 0; i < table.size; i++ {
		key := string([]rune{'A', rune('A' + i)})
		if err := table.Insert(key, "X"); err != nil {
			t.Fatalf("insert failed at %d: %v", i, err)
		}
	}
	if err := table.Insert("ZZ", "X"); err != ErrTableFull {
		t.Fatalf("expected ErrTableFull, got %v", err)
	}
}

func TestRenderOutput(t *testing.T) {
	table := newTestTable(t)
	if err := table.Insert("Анатомия", "A"); err != nil {
		t.Fatalf("insert failed: %v", err)
	}
	output := table.Render()
	if output == "" {
		t.Fatalf("expected non-empty render output")
	}
}

func newTestTable(t *testing.T) *Table {
	table, err := NewTable(DefaultTableSize)
	if err != nil {
		t.Fatalf("new table failed: %v", err)
	}
	return table
}

func stringsUpper(value string) string {
	upper, err := normalizeKey(value)
	if err != nil {
		return value
	}
	return upper
}

func assertFloatNear(t *testing.T, got float64, expected float64) {
	t.Helper()
	const epsilon = 0.0001
	if math.Abs(got-expected) > epsilon {
		t.Fatalf("expected %.2f, got %.2f", expected, got)
	}
}
