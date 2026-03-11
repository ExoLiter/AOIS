package minimization

import (
	"fmt"
	"logical_calculator/internal/models"
	"sort"
	"strings"
)

func BuildMinterms(table []models.Row, vars []string) []models.Term {
	var terms []models.Term
	for _, row := range table {
		if row.Result {
			terms = append(terms, createBaseTerm(row.Index, vars))
		}
	}
	return terms
}

func BuildMaxterms(table []models.Row, vars []string) []models.Term {
	var terms []models.Term
	for _, row := range table {
		if !row.Result {
			terms = append(terms, createBaseTerm(row.Index, vars))
		}
	}
	return terms
}

func createBaseTerm(index int, vars []string) models.Term {
	return models.Term{
		Mask:    (1 << len(vars)) - 1,
		Value:   index,
		Indices: []int{index},
	}
}

func TryMerge(t1, t2 models.Term) (models.Term, bool) {
	if t1.Mask != t2.Mask {
		return models.Term{}, false
	}
	diff := t1.Value ^ t2.Value
	if diff != 0 && (diff&(diff-1)) == 0 {
		newMask := t1.Mask &^ diff
		newValue := t1.Value &^ diff
		newIndices := append(append([]int{}, t1.Indices...), t2.Indices...)
		sort.Ints(newIndices)
		return models.Term{Mask: newMask, Value: newValue, Indices: newIndices}, true
	}
	return models.Term{}, false
}

// FormatTermSDNF Форматирует для СДНФ (a & !b)
func FormatTermSDNF(t models.Term, vars []string) string {
	var parts []string
	for i, v := range vars {
		shift := len(vars) - 1 - i
		if (t.Mask>>shift)&1 == 1 {
			if (t.Value>>shift)&1 == 1 {
				parts = append(parts, v)
			} else {
				parts = append(parts, "!"+v)
			}
		}
	}
	if len(parts) == 0 {
		return "1"
	}
	return formatWithIndices(strings.Join(parts, ""), t.Indices)
}

// FormatTermSKNF Форматирует для СКНФ (!a v b)
func FormatTermSKNF(t models.Term, vars []string) string {
	var parts []string
	for i, v := range vars {
		shift := len(vars) - 1 - i
		if (t.Mask>>shift)&1 == 1 {
			if (t.Value>>shift)&1 == 1 {
				parts = append(parts, "!"+v) // В СКНФ 1 - это отрицание
			} else {
				parts = append(parts, v)
			}
		}
	}
	if len(parts) == 0 {
		return "0"
	}
	return formatWithIndices(strings.Join(parts, " v "), t.Indices)
}

func formatWithIndices(expr string, indices []int) string {
	idxStrs := make([]string, len(indices))
	for i, idx := range indices {
		idxStrs[i] = fmt.Sprint(idx)
	}
	return fmt.Sprintf("(%s)%s", expr, strings.Join(idxStrs, ","))
}

func FormatTermsSum(terms []models.Term, vars []string) string {
	if len(terms) == 0 {
		return "0"
	}
	var parts []string
	for _, t := range terms {
		parts = append(parts, FormatTermSDNF(t, vars))
	}
	return strings.Join(parts, " v ")
}

func FormatTermsProd(terms []models.Term, vars []string) string {
	if len(terms) == 0 {
		return "1"
	}
	var parts []string
	for _, t := range terms {
		parts = append(parts, FormatTermSKNF(t, vars))
	}
	return strings.Join(parts, " & ")
}
