package truthtable

import (
	"logical_calculator/internal/evaluator"
	"logical_calculator/internal/models"
	"sort"
)

// ExtractVariables находит уникальные переменные
func ExtractVariables(rpn []string) []string {
	varSet := make(map[string]bool)
	for _, token := range rpn {
		if token >= "a" && token <= "e" {
			varSet[token] = true
		}
	}
	var vars []string
	for v := range varSet {
		vars = append(vars, v)
	}
	sort.Strings(vars)
	return vars
}

// GenerateTable строит таблицу истинности (до 2^5 = 32 строк)
func GenerateTable(rpn []string, vars []string) []models.Row {
	var table []models.Row
	rowCount := 1 << len(vars)

	for i := 0; i < rowCount; i++ {
		values := createRowValues(i, vars)
		res := evaluator.EvaluateRPN(rpn, values)
		table = append(table, models.Row{
			Values: values,
			Result: res,
			Index:  i,
		})
	}
	return table
}

func createRowValues(index int, vars []string) map[string]bool {
	values := make(map[string]bool)
	for bitIndex, v := range vars {
		shift := len(vars) - 1 - bitIndex
		values[v] = ((index >> shift) & 1) == 1
	}
	return values
}
