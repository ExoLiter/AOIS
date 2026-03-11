package zhegalkin

import (
	"logical_calculator/internal/models"
	"strings"
)

func BuildPolynomial(table []models.Row, vars []string) string {
	coeffs := ComputeCoeffs(table)
	var terms []string

	for i, coeff := range coeffs {
		if coeff {
			terms = append(terms, buildZhegalkinTerm(i, vars))
		}
	}
	if len(terms) == 0 {
		return "0"
	}
	return strings.Join(terms, " + ")
}

// ComputeCoeffs вычисляет коэффициенты полинома (сделана публичной)
func ComputeCoeffs(table []models.Row) []bool {
	n := len(table)
	coeffs := make([]bool, n)
	for i, row := range table {
		coeffs[i] = row.Result
	}
	for i := 1; i < n; i *= 2 {
		for j := 0; j < n; j += 2 * i {
			for k := 0; k < i; k++ {
				coeffs[j+k+i] = coeffs[j+k+i] != coeffs[j+k]
			}
		}
	}
	return coeffs
}

func buildZhegalkinTerm(index int, vars []string) string {
	if index == 0 {
		return "1"
	}
	var parts []string
	for bitIndex, v := range vars {
		shift := len(vars) - 1 - bitIndex
		if (index>>shift)&1 == 1 {
			parts = append(parts, v)
		}
	}
	return strings.Join(parts, "")
}
