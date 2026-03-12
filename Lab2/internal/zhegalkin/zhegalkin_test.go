package zhegalkin

import (
	"logical_calculator/internal/models"
	"testing"
)

func TestComputeCoeffsXOR(t *testing.T) {
	// XOR truth table for a,b: 00->0, 01->1, 10->1, 11->0
	table := []models.Row{
		{Index: 0, Result: false},
		{Index: 1, Result: true},
		{Index: 2, Result: true},
		{Index: 3, Result: false},
	}
	coeffs := ComputeCoeffs(table)
	want := []bool{false, true, true, false}
	for i := range want {
		if coeffs[i] != want[i] {
			t.Fatalf("coeff[%d] mismatch: got %v want %v", i, coeffs[i], want[i])
		}
	}
}

func TestBuildPolynomial(t *testing.T) {
	zeroTable := []models.Row{
		{Index: 0, Result: false},
		{Index: 1, Result: false},
		{Index: 2, Result: false},
		{Index: 3, Result: false},
	}
	if got := BuildPolynomial(zeroTable, []string{"a", "b"}); got != "0" {
		t.Fatalf("expected zero polynomial, got %q", got)
	}

	xorTable := []models.Row{
		{Index: 0, Result: false},
		{Index: 1, Result: true},
		{Index: 2, Result: true},
		{Index: 3, Result: false},
	}
	if got := BuildPolynomial(xorTable, []string{"a", "b"}); got != "b + a" {
		t.Fatalf("unexpected polynomial: %q", got)
	}
}
