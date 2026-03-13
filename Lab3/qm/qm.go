package qm

import (
	"math/bits"
	"strings"
)

// Implicant представляет собой импликанту функции (маска показывает "Don't care" биты).
type Implicant struct {
	Value int
	Mask  int
}

// IsEqual проверяет, идентичны ли две импликанты.
func (a Implicant) IsEqual(b Implicant) bool {
	return a.Value == b.Value && a.Mask == b.Mask
}

// Covers проверяет, покрывает ли импликанта заданный минтерм.
func (imp Implicant) Covers(minterm int) bool {
	return (minterm & ^imp.Mask) == (imp.Value & ^imp.Mask)
}

// DifferByOneBit проверяет возможность склейки двух импликант.
func DifferByOneBit(a, b Implicant) (bool, Implicant) {
	if a.Mask != b.Mask {
		return false, Implicant{}
	}
	diff := uint(a.Value ^ b.Value)
	if bits.OnesCount(diff) == 1 {
		return true, Implicant{Value: a.Value &^ int(diff), Mask: a.Mask | int(diff)}
	}
	return false, Implicant{}
}

// GenerateSDNF генерирует СДНФ для заданных минтермов.
func GenerateSDNF(numVars int, minterms []int, varNames []string) string {
	if len(minterms) == 0 {
		return "0"
	}
	var parts []string
	for _, m := range minterms {
		parts = append(parts, formatImplicant(Implicant{Value: m, Mask: 0}, numVars, varNames))
	}
	return strings.Join(parts, " | ")
}

// Minimize выполняет минимизацию методом Куайна-Маккласки.
func Minimize(numVars int, minterms []int, dontCares []int, varNames []string) string {
	if len(minterms) == 0 {
		return "0"
	}
	primes := findPrimeImplicants(minterms, dontCares)
	essentials, remaining := findEssentialPrimes(primes, minterms)
	solution := coverRemaining(remaining, primes, essentials)
	return formatSolution(solution, numVars, varNames)
}

func findPrimeImplicants(minterms, dontCares []int) []Implicant {
	implicants := initImplicants(minterms, dontCares)
	primes := make(map[Implicant]bool)

	for len(implicants) > 0 {
		nextLevel := make(map[Implicant]bool)
		used := make(map[Implicant]bool)

		for i, a := range implicants {
			for j := i + 1; j < len(implicants); j++ {
				b := implicants[j]
				if canMerge, merged := DifferByOneBit(a, b); canMerge {
					nextLevel[merged] = true
					used[a] = true
					used[b] = true
				}
			}
		}

		for _, imp := range implicants {
			if !used[imp] {
				primes[imp] = true
			}
		}

		implicants = make([]Implicant, 0, len(nextLevel))
		for imp := range nextLevel {
			implicants = append(implicants, imp)
		}
	}

	result := make([]Implicant, 0, len(primes))
	for imp := range primes {
		result = append(result, imp)
	}
	return result
}

func initImplicants(minterms, dontCares []int) []Implicant {
	var res []Implicant
	for _, m := range minterms {
		res = append(res, Implicant{Value: m, Mask: 0})
	}
	for _, m := range dontCares {
		res = append(res, Implicant{Value: m, Mask: 0})
	}
	return res
}

func findEssentialPrimes(primes []Implicant, minterms []int) ([]Implicant, []int) {
	var essentials []Implicant
	covered := make(map[int]bool)

	for _, m := range minterms {
		covers := getCovers(primes, m)
		if len(covers) == 1 {
			essentials = appendUnique(essentials, covers[0])
		}
	}

	for _, e := range essentials {
		for _, m := range minterms {
			if e.Covers(m) {
				covered[m] = true
			}
		}
	}

	var remaining []int
	for _, m := range minterms {
		if !covered[m] {
			remaining = append(remaining, m)
		}
	}
	return essentials, remaining
}

func getCovers(primes []Implicant, minterm int) []Implicant {
	var res []Implicant
	for _, p := range primes {
		if p.Covers(minterm) {
			res = append(res, p)
		}
	}
	return res
}

func appendUnique(list []Implicant, item Implicant) []Implicant {
	for _, x := range list {
		if x.IsEqual(item) {
			return list
		}
	}
	return append(list, item)
}

func coverRemaining(remaining []int, primes, essentials []Implicant) []Implicant {
	solution := append([]Implicant{}, essentials...)
	uncovered := append([]int{}, remaining...)

	for len(uncovered) > 0 {
		best := findBestPrime(primes, uncovered)
		solution = append(solution, best)
		var nextUncovered []int
		for _, m := range uncovered {
			if !best.Covers(m) {
				nextUncovered = append(nextUncovered, m)
			}
		}
		uncovered = nextUncovered
	}
	return solution
}

func findBestPrime(primes []Implicant, uncovered []int) Implicant {
	bestCount := -1
	var best Implicant
	for _, p := range primes {
		count := 0
		for _, m := range uncovered {
			if p.Covers(m) {
				count++
			}
		}
		if count > bestCount {
			bestCount = count
			best = p
		}
	}
	return best
}

func formatSolution(solution []Implicant, numVars int, varNames []string) string {
	var parts []string
	for _, imp := range solution {
		if imp.Mask == (1<<numVars)-1 {
			return "1"
		}
		parts = append(parts, formatImplicant(imp, numVars, varNames))
	}
	return strings.Join(parts, " | ")
}

func formatImplicant(imp Implicant, numVars int, varNames []string) string {
	var parts []string
	for i := 0; i < numVars; i++ {
		bitPos := numVars - 1 - i
		if ((imp.Mask >> bitPos) & 1) == 0 {
			if ((imp.Value >> bitPos) & 1) == 1 {
				parts = append(parts, varNames[i])
			} else {
				parts = append(parts, "!"+varNames[i])
			}
		}
	}
	return "(" + strings.Join(parts, " & ") + ")"
}
