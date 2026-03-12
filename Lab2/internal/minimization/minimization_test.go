package minimization

import (
	"logical_calculator/internal/models"
	"logical_calculator/internal/truthtable"
	"strings"
	"testing"
)

func TestTryMerge(t *testing.T) {
	t1 := models.Term{Mask: 3, Value: 0, Indices: []int{0}}
	t2 := models.Term{Mask: 3, Value: 1, Indices: []int{1}}
	merged, ok := TryMerge(t1, t2)
	if !ok {
		t.Fatalf("expected merge to succeed")
	}
	if merged.Mask != 2 || merged.Value != 0 || len(merged.Indices) != 2 {
		t.Fatalf("unexpected merged term: %#v", merged)
	}

	t3 := models.Term{Mask: 3, Value: 3, Indices: []int{3}}
	if _, ok := TryMerge(t1, t3); ok {
		t.Fatalf("expected merge to fail for diff >1")
	}
}

func TestFormatTermsEdgeCases(t *testing.T) {
	term := models.Term{Mask: 0, Value: 0, Indices: []int{0}}
	if got := FormatTermSDNF(term, []string{"a", "b"}); got != "1" {
		t.Fatalf("unexpected SDNF term: %q", got)
	}
	if got := FormatTermSKNF(term, []string{"a", "b"}); got != "0" {
		t.Fatalf("unexpected SKNF term: %q", got)
	}

	if got := FormatTermsSum(nil, []string{"a"}); got != "0" {
		t.Fatalf("unexpected empty sum: %q", got)
	}
	if got := FormatTermsProd(nil, []string{"a"}); got != "1" {
		t.Fatalf("unexpected empty product: %q", got)
	}
}

func TestMinimizeCalculationAndTabular(t *testing.T) {
	vars := []string{"a", "b"}
	table := truthtable.GenerateTable([]string{"a", "b", "|"}, vars)

	outCalc := MinimizeCalculationSDNF(table, vars)
	if !strings.Contains(outCalc, "(a)") || !strings.Contains(outCalc, "(b)") {
		t.Fatalf("calculation output missing expected terms: %q", outCalc)
	}

	outTab := MinimizeTabularCalcSDNF(table, vars)
	if !strings.Contains(outTab, " X") {
		t.Fatalf("tabular output missing table marks: %q", outTab)
	}
}

func TestMinimizeKarnaughAndEmpty(t *testing.T) {
	vars := []string{"a", "b"}
	table := truthtable.GenerateTable([]string{"a", "b", "&"}, vars)
	outK := MinimizeKarnaughSDNF(table, vars)
	if !strings.Contains(outK, "0") && !strings.Contains(outK, "1") {
		t.Fatalf("karnaugh output looks empty: %q", outK)
	}

	// Empty SDNF (all false) should short-circuit.
	allFalse := []models.Row{
		{Index: 0, Result: false},
		{Index: 1, Result: false},
	}
	outEmpty := MinimizeCalculationSDNF(allFalse, []string{"a"})
	if !strings.Contains(outEmpty, "Нет импликант") {
		t.Fatalf("unexpected empty output: %q", outEmpty)
	}
}

func TestBuildMaxtermsAndSKNFEmpty(t *testing.T) {
	vars := []string{"a"}
	table := []models.Row{
		{Index: 0, Result: true},
		{Index: 1, Result: true},
	}
	maxterms := BuildMaxterms(table, vars)
	if len(maxterms) != 0 {
		t.Fatalf("expected no maxterms, got %#v", maxterms)
	}

	outCalc := MinimizeCalculationSKNF(table, vars)
	if !strings.Contains(outCalc, "Нет импликант") {
		t.Fatalf("unexpected empty SKNF output: %q", outCalc)
	}

	outTab := MinimizeTabularCalcSKNF(table, vars)
	if !strings.Contains(outTab, "Нет импликант") {
		t.Fatalf("unexpected empty tabular SKNF output: %q", outTab)
	}
}

func TestBuildQuineTableAndHelpers(t *testing.T) {
	vars := []string{"a", "b"}
	terms := []models.Term{
		{Mask: 3, Value: 1, Indices: []int{1}},
		{Mask: 3, Value: 3, Indices: []int{3}},
	}
	table := buildQuineTable(terms, terms, vars, FormatTermSDNF)
	if !strings.Contains(table, "X") {
		t.Fatalf("expected table marks, got: %q", table)
	}

	if got := formatBinary(3, 2); got != "11" {
		t.Fatalf("unexpected binary format: %q", got)
	}

	rows, cols := splitVars([]string{"a", "b", "c"})
	if len(rows) != 1 || len(cols) != 2 {
		t.Fatalf("unexpected split vars: %v %v", rows, cols)
	}

	if codes := generateGrayCodes(2); len(codes) != 4 || codes[1] != 1 {
		t.Fatalf("unexpected gray codes: %#v", codes)
	}
}

func TestBuildKarnaughMapLimits(t *testing.T) {
	vars := []string{"a", "b", "c", "d", "e"}
	table := truthtable.GenerateTable([]string{"a"}, vars)
	out := MinimizeKarnaughSDNF(table, vars)
	if !strings.Contains(out, ">4") {
		t.Fatalf("expected >4 warning, got: %q", out)
	}
}
