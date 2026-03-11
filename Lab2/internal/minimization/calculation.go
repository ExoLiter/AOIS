package minimization

import (
	"fmt"
	"logical_calculator/internal/models"
	"strings"
)

func MinimizeCalculationSDNF(table []models.Row, vars []string) string {
	minterms := BuildMinterms(table, vars)
	return runCalculationMethod(minterms, vars, FormatTermSDNF, FormatTermsSum, "v")
}

func MinimizeCalculationSKNF(table []models.Row, vars []string) string {
	maxterms := BuildMaxterms(table, vars)
	return runCalculationMethod(maxterms, vars, FormatTermSKNF, FormatTermsProd, "&")
}

// Внутренняя универсальная функция для СДНФ и СКНФ
func runCalculationMethod(terms []models.Term, vars []string,
	formatTerm func(models.Term, []string) string,
	formatFinal func([]models.Term, []string) string, op string) string {

	if len(terms) == 0 {
		return "Нет импликант для склеивания."
	}

	var sb strings.Builder
	sb.WriteString("Этап склеивания:\n")

	primeImplicants := PerformGluingSteps(&sb, terms, vars, formatTerm, op)

	sb.WriteString("Результат:\n" + formatFinal(primeImplicants, vars) + "\n\n")
	sb.WriteString("Проверка на лишние импликанты:\n")

	finalImplicants := RemoveRedundantCalc(&sb, primeImplicants, terms, vars, formatTerm)
	sb.WriteString("\nИтоговый результат:\n" + formatFinal(finalImplicants, vars) + "\n")

	return sb.String()
}

func PerformGluingSteps(sb *strings.Builder, current []models.Term, vars []string,
	formatTerm func(models.Term, []string) string, op string) []models.Term {
	var primes []models.Term
	for {
		next, used, mergedStr := glueRound(current, vars, formatTerm, op)
		if mergedStr != "" && sb != nil {
			sb.WriteString(mergedStr)
		}
		primes = extractUnused(current, used, primes)
		if len(next) == 0 {
			break
		}
		current = next
	}
	return removeDuplicates(primes)
}

func glueRound(current []models.Term, vars []string, formatTerm func(models.Term, []string) string, op string) ([]models.Term, map[int]bool, string) {
	var next []models.Term
	used := make(map[int]bool)
	var sb strings.Builder

	for i := 0; i < len(current); i++ {
		for j := i + 1; j < len(current); j++ {
			merged, ok := TryMerge(current[i], current[j])
			if ok {
				used[i], used[j] = true, true
				next = append(next, merged)
				sb.WriteString(fmt.Sprintf("%s %s %s => %s\n",
					formatTerm(current[i], vars), op, formatTerm(current[j], vars), formatTerm(merged, vars)))
			}
		}
	}
	return next, used, sb.String()
}

func extractUnused(current []models.Term, used map[int]bool, primes []models.Term) []models.Term {
	for i, t := range current {
		if !used[i] {
			primes = append(primes, t)
		}
	}
	return primes
}

func removeDuplicates(terms []models.Term) []models.Term {
	seen := make(map[int]bool)
	var res []models.Term
	for _, t := range terms {
		hash := (t.Mask << 16) | t.Value
		if !seen[hash] {
			seen[hash] = true
			res = append(res, t)
		}
	}
	return res
}

func RemoveRedundantCalc(sb *strings.Builder, primes []models.Term, baseTerms []models.Term, vars []string, formatTerm func(models.Term, []string) string) []models.Term {
	var final []models.Term
	for i, candidate := range primes {
		if isEssential(candidate, primes, baseTerms, i) {
			final = append(final, candidate)
			if sb != nil {
				sb.WriteString(fmt.Sprintf("%s - эта импликанта не лишняя\n", formatTerm(candidate, vars)))
			}
		} else {
			if sb != nil {
				sb.WriteString(fmt.Sprintf("%s - эта импликанта лишняя\n", formatTerm(candidate, vars)))
			}
		}
	}
	return final
}

func isEssential(target models.Term, all []models.Term, baseTerms []models.Term, skipIdx int) bool {
	for _, m := range baseTerms {
		if covers(target, m) {
			coveredByOthers := false
			for j, other := range all {
				if j != skipIdx && covers(other, m) {
					coveredByOthers = true
					break
				}
			}
			if !coveredByOthers {
				return true
			}
		}
	}
	return false
}

func covers(implicant, minterm models.Term) bool {
	return (minterm.Value & implicant.Mask) == implicant.Value
}
