package forms

import (
	"logical_calculator/internal/models"
	"strings"
)

func BuildSDNF(table []models.Row, vars []string) string {
	var terms []string
	for _, row := range table {
		if row.Result {
			terms = append(terms, buildTerm(row, vars, true))
		}
	}
	return strings.Join(terms, " v ")
}

func BuildSKNF(table []models.Row, vars []string) string {
	var terms []string
	for _, row := range table {
		if !row.Result {
			terms = append(terms, buildTerm(row, vars, false))
		}
	}
	return strings.Join(terms, " & ")
}

func buildTerm(row models.Row, vars []string, isSdnf bool) string {
	var parts []string
	for _, v := range vars {
		parts = append(parts, formatVar(v, row.Values[v], isSdnf))
	}
	separator := " & "
	if !isSdnf {
		separator = " v "
	}
	return "(" + strings.Join(parts, separator) + ")"
}

func formatVar(v string, val bool, isSdnf bool) string {
	needsNot := (!val && isSdnf) || (val && !isSdnf)
	if needsNot {
		return "!" + v
	}
	return v
}

func NumericForms(table []models.Row) (sdnfNums, sknfNums []int) {
	for _, row := range table {
		if row.Result {
			sdnfNums = append(sdnfNums, row.Index)
		} else {
			sknfNums = append(sknfNums, row.Index)
		}
	}
	return sdnfNums, sknfNums
}

func IndexForm(table []models.Row) int {
	index := 0
	for i, row := range table {
		if row.Result {
			shift := len(table) - 1 - i
			index |= (1 << shift)
		}
	}
	return index
}
