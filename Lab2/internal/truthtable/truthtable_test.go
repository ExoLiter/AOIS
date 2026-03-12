package truthtable

import (
	"reflect"
	"testing"
)

func TestExtractVariablesSorted(t *testing.T) {
	rpn := []string{"b", "a", "&", "c", "|"}
	got := ExtractVariables(rpn)
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ExtractVariables mismatch: got %#v want %#v", got, want)
	}
}

func TestGenerateTableAndValues(t *testing.T) {
	rpn := []string{"a", "b", "&"}
	vars := []string{"a", "b"}
	table := GenerateTable(rpn, vars)
	if len(table) != 4 {
		t.Fatalf("expected 4 rows, got %d", len(table))
	}

	if table[0].Index != 0 || table[0].Values["a"] || table[0].Values["b"] || table[0].Result {
		t.Fatalf("row 0 mismatch: %#v", table[0])
	}
	if table[3].Index != 3 || !table[3].Values["a"] || !table[3].Values["b"] || !table[3].Result {
		t.Fatalf("row 3 mismatch: %#v", table[3])
	}
}
