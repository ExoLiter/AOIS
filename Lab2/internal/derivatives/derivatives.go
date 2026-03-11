package derivatives

import (
	"fmt"
	"logical_calculator/internal/minimization"
	"logical_calculator/internal/models"
	"strings"
)

func FindDummyVariables(table []models.Row, vars []string) []string {
	var dummies []string
	for i, v := range vars {
		if isDummy(table, len(vars), i) {
			dummies = append(dummies, v)
		}
	}
	return dummies
}

func isDummy(table []models.Row, varCount, targetVarIndex int) bool {
	shift := varCount - 1 - targetVarIndex
	for i := 0; i < len(table); i++ {
		if (i>>shift)&1 == 0 {
			flippedIdx := i | (1 << shift)
			if table[i].Result != table[flippedIdx].Result {
				return false
			}
		}
	}
	return true
}

// GenerateAllDerivatives собирает полный отчет по всем производным
func GenerateAllDerivatives(table []models.Row, vars []string) string {
	var sb strings.Builder
	baseResults := make([]bool, len(table))
	for i, r := range table {
		baseResults[i] = r.Result
	}

	for order := 1; order <= len(vars) && order <= 4; order++ {
		sb.WriteString(fmt.Sprintf("\n--- Производные %d-го порядка ---\n", order))
		combinations := getCombinations(vars, order)
		for _, comb := range combinations {
			processDerivativeCombo(&sb, baseResults, vars, comb)
		}
	}
	return sb.String()
}

func processDerivativeCombo(sb *strings.Builder, baseResults []bool, vars []string, comb []string) {
	currentRes := append([]bool(nil), baseResults...) // Копия базового массива
	for _, tVar := range comb {
		idx := findVarIndex(vars, tVar)
		currentRes = applyDerivative(currentRes, len(vars), idx)
	}

	indicesStr := getIndicesStr(currentRes)
	minimizedStr := minimizeDerivativeResult(currentRes, vars)

	title := "d"
	if len(comb) > 1 {
		title = fmt.Sprintf("d^%d", len(comb))
	}
	sb.WriteString(fmt.Sprintf("%sF / d%s:\n", title, strings.Join(comb, " d")))
	sb.WriteString(fmt.Sprintf("  Индексы: %s\n", indicesStr))
	sb.WriteString(fmt.Sprintf("  Функция: %s\n", minimizedStr))
}

func applyDerivative(results []bool, varCount, targetIndex int) []bool {
	deriv := make([]bool, len(results))
	shift := varCount - 1 - targetIndex
	for i := 0; i < len(results); i++ {
		flippedIdx := i ^ (1 << shift)
		deriv[i] = results[i] != results[flippedIdx]
	}
	return deriv
}

// minimizeDerivativeResult превращает вектор производной в СДНФ и упрощает её
func minimizeDerivativeResult(results []bool, vars []string) string {
	// Создаем фейковую таблицу для минимизатора
	var fakeTable []models.Row
	for i, res := range results {
		fakeTable = append(fakeTable, models.Row{Index: i, Result: res})
	}

	minterms := minimization.BuildMinterms(fakeTable, vars)
	if len(minterms) == 0 {
		return "0"
	}

	primes := minimization.PerformGluingSteps(nil, minterms, vars, minimization.FormatTermSDNF, "v")
	final := minimization.RemoveRedundantCalc(nil, primes, minterms, vars, minimization.FormatTermSDNF)

	return minimization.FormatTermsSum(final, vars)
}

func findVarIndex(vars []string, target string) int {
	for i, v := range vars {
		if v == target {
			return i
		}
	}
	return -1
}

func getIndicesStr(results []bool) string {
	var indices []string
	for i, res := range results {
		if res {
			indices = append(indices, fmt.Sprint(i))
		}
	}
	if len(indices) == 0 {
		return "Нет наборов (Всегда 0)"
	}
	return "[" + strings.Join(indices, ", ") + "]"
}

// getCombinations генерирует все уникальные сочетания переменных заданного размера
func getCombinations(vars []string, r int) [][]string {
	var result [][]string
	var generate func(start int, current []string)
	generate = func(start int, current []string) {
		if len(current) == r {
			comb := make([]string, len(current))
			copy(comb, current)
			result = append(result, comb)
			return
		}
		for i := start; i < len(vars); i++ {
			generate(i+1, append(current, vars[i]))
		}
	}
	generate(0, []string{})
	return result
}
