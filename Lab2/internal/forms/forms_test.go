package forms

import (
	"logical_calculator/internal/truthtable"
	"testing"
)

func TestBuildFormsAndNumericIndex(t *testing.T) {
	rpn := []string{"a", "b", "|"}
	vars := []string{"a", "b"}
	table := truthtable.GenerateTable(rpn, vars)

	sdnf := BuildSDNF(table, vars)
	if sdnf != "(!a & b) v (a & !b) v (a & b)" {
		t.Fatalf("unexpected SDNF: %q", sdnf)
	}

	sknf := BuildSKNF(table, vars)
	if sknf != "(a v b)" {
		t.Fatalf("unexpected SKNF: %q", sknf)
	}

	sdnfNums, sknfNums := NumericForms(table)
	if len(sdnfNums) != 3 || sdnfNums[0] != 1 || sdnfNums[1] != 2 || sdnfNums[2] != 3 {
		t.Fatalf("unexpected SDNF numbers: %#v", sdnfNums)
	}
	if len(sknfNums) != 1 || sknfNums[0] != 0 {
		t.Fatalf("unexpected SKNF numbers: %#v", sknfNums)
	}

	if idx := IndexForm(table); idx != 7 {
		t.Fatalf("unexpected index form: %d", idx)
	}
}
