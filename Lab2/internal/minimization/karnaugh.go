package minimization

import (
	"fmt"
	"logical_calculator/internal/models"
	"strings"
)

func MinimizeKarnaughSDNF(table []models.Row, vars []string) string {
	minterms := BuildMinterms(table, vars)
	finalFunc := getSilentMinimized(minterms, vars, FormatTermSDNF)
	// Убрали передачу параметра isSdnf
	return buildKarnaughMap(table, vars, FormatTermsSum(finalFunc, vars))
}

func MinimizeKarnaughSKNF(table []models.Row, vars []string) string {
	maxterms := BuildMaxterms(table, vars)
	finalFunc := getSilentMinimized(maxterms, vars, FormatTermSKNF)
	// Убрали передачу параметра isSdnf
	return buildKarnaughMap(table, vars, FormatTermsProd(finalFunc, vars))
}

func getSilentMinimized(terms []models.Term, vars []string, formatTerm func(models.Term, []string) string) []models.Term {
	primes := PerformGluingSteps(nil, terms, vars, formatTerm, "")
	return RemoveRedundantCalc(nil, primes, terms, vars, formatTerm)
}

// Изменили сигнатуру функции (удалили isSdnf bool)
func buildKarnaughMap(table []models.Row, vars []string, finalResult string) string {
	if len(vars) > 4 {
		return "Отрисовка карты Карно для >4 переменных выходит за рамки текстовой консоли.\n"
	}

	rowVars, colVars := splitVars(vars)
	rowCodes := generateGrayCodes(len(rowVars))
	colCodes := generateGrayCodes(len(colVars))

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s \\ %s\t|", strings.Join(rowVars, ""), strings.Join(colVars, "")))
	for _, c := range colCodes {
		sb.WriteString(fmt.Sprintf(" %s |", formatBinary(c, len(colVars))))
	}
	sb.WriteString("\n" + strings.Repeat("-", 10+5*len(colCodes)) + "\n")

	for _, r := range rowCodes {
		sb.WriteString(fmt.Sprintf("  %s\t|", formatBinary(r, len(rowVars))))
		for _, c := range colCodes {
			idx := (r << len(colVars)) | c
			res := getTableResult(table, idx)

			val := "0"
			if res {
				val = "1"
			}
			sb.WriteString(fmt.Sprintf("  %s |", val))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\nРезультат минимизации по карте Карно:\n" + finalResult + "\n")
	return sb.String()
}

func splitVars(vars []string) ([]string, []string) {
	half := len(vars) / 2
	if len(vars) == 3 {
		half = 1
	}
	return vars[:half], vars[half:]
}

func generateGrayCodes(bits int) []int {
	count := 1 << bits
	codes := make([]int, count)
	for i := 0; i < count; i++ {
		codes[i] = i ^ (i >> 1)
	}
	return codes
}

func formatBinary(val, bits int) string {
	format := fmt.Sprintf("%%0%db", bits)
	return fmt.Sprintf(format, val)
}

func getTableResult(table []models.Row, index int) bool {
	for _, r := range table {
		if r.Index == index {
			return r.Result
		}
	}
	return false
}
