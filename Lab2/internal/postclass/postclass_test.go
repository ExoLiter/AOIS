package postclass

import (
	"logical_calculator/internal/truthtable"
	"testing"
)

func TestPostClasses(t *testing.T) {
	vars2 := []string{"a", "b"}

	orTable := truthtable.GenerateTable([]string{"a", "b", "|"}, vars2)
	if !IsT0(orTable) || !IsT1(orTable) {
		t.Fatalf("OR should be in T0 and T1")
	}

	notTable := truthtable.GenerateTable([]string{"a", "!"}, []string{"a"})
	if !IsSelfDual(notTable) {
		t.Fatalf("NOT a should be self-dual")
	}

	andTable := truthtable.GenerateTable([]string{"a", "b", "&"}, vars2)
	if !IsMonotonic(andTable) {
		t.Fatalf("AND should be monotonic")
	}

	xorTable := truthtable.GenerateTable([]string{"a", "b", "=", "!"}, vars2)
	if IsMonotonic(xorTable) {
		t.Fatalf("XOR should not be monotonic")
	}
	if !IsLinear(xorTable) {
		t.Fatalf("XOR should be linear")
	}
	if IsLinear(andTable) {
		t.Fatalf("AND should not be linear")
	}

	nandTable := truthtable.GenerateTable([]string{"a", "b", "&", "!"}, vars2)
	if !IsFunctionallyComplete(nandTable) {
		t.Fatalf("NAND should be functionally complete")
	}
}
