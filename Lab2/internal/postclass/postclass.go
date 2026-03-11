package postclass

import (
	"logical_calculator/internal/models"
	"logical_calculator/internal/zhegalkin"
)

func IsT0(table []models.Row) bool {
	return !table[0].Result
}

func IsT1(table []models.Row) bool {
	return table[len(table)-1].Result
}

func IsSelfDual(table []models.Row) bool {
	n := len(table)
	for i := 0; i < n/2; i++ {
		if table[i].Result == table[n-1-i].Result {
			return false
		}
	}
	return true
}

func IsMonotonic(table []models.Row) bool {
	for i := 0; i < len(table); i++ {
		for j := i + 1; j < len(table); j++ {
			if isSubset(i, j) && table[i].Result && !table[j].Result {
				return false
			}
		}
	}
	return true
}

// IsLinear проверяет класс L. Функция линейна, если в полиноме нет конъюнкций (умножений)
func IsLinear(table []models.Row) bool {
	coeffs := zhegalkin.ComputeCoeffs(table)
	for i, c := range coeffs {
		if c && countBits(i) > 1 {
			return false // Нашли коэффициент = 1 для члена с >1 переменными
		}
	}
	return true
}

// IsFunctionallyComplete проверяет, является ли функция функционально полной по теореме Поста
func IsFunctionallyComplete(table []models.Row) bool {
	return !IsT0(table) && !IsT1(table) && !IsSelfDual(table) && !IsMonotonic(table) && !IsLinear(table)
}

func isSubset(a, b int) bool {
	return (a & b) == a
}

func countBits(n int) int {
	count := 0
	for n > 0 {
		count += n & 1
		n >>= 1
	}
	return count
}
