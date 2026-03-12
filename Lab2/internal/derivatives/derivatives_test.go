package derivatives

import (
	"logical_calculator/internal/truthtable"
	"strings"
	"testing"
)

func TestFindDummyVariables(t *testing.T) {
	rpn := []string{"a"}
	vars := []string{"a", "b"}
	table := truthtable.GenerateTable(rpn, vars)
	dummies := FindDummyVariables(table, vars)
	if len(dummies) != 1 || dummies[0] != "b" {
		t.Fatalf("unexpected dummies: %#v", dummies)
	}
}

func TestApplyDerivative(t *testing.T) {
	// AND for a,b: 00->0, 01->0, 10->0, 11->1
	results := []bool{false, false, false, true}
	deriv := applyDerivative(results, 2, 0) // d/da
	want := []bool{false, true, false, true}
	for i := range want {
		if deriv[i] != want[i] {
			t.Fatalf("derivative[%d] mismatch: got %v want %v", i, deriv[i], want[i])
		}
	}
}

func TestGenerateAllDerivativesOutput(t *testing.T) {
	rpn := []string{"a", "b", "|"}
	vars := []string{"a", "b"}
	table := truthtable.GenerateTable(rpn, vars)

	out := GenerateAllDerivatives(table, vars)
	if !strings.Contains(out, "dF / da:") {
		t.Fatalf("output missing first order derivative header")
	}
}
