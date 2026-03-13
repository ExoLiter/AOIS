package qm

import "testing"

func TestImplicantIsEqual(t *testing.T) {
	a := Implicant{Value: 1, Mask: 2}
	b := Implicant{Value: 1, Mask: 2}
	c := Implicant{Value: 2, Mask: 2}
	if !a.IsEqual(b) {
		t.Error("IsEqual failed")
	}
	if a.IsEqual(c) {
		t.Error("IsEqual false positive")
	}
}

func TestDifferByOneBit(t *testing.T) {
	a := Implicant{Value: 0, Mask: 0}
	b := Implicant{Value: 1, Mask: 0}
	ok, m := DifferByOneBit(a, b)
	if !ok || m.Value != 0 || m.Mask != 1 {
		t.Error("DifferByOneBit failed")
	}

	// Same mask, but differences in multiple bits.
	c := Implicant{Value: 3, Mask: 0}
	if ok, _ = DifferByOneBit(a, c); ok {
		t.Error("DifferByOneBit false positive on multiple bits")
	}

	// Different masks should not merge.
	d := Implicant{Value: 0, Mask: 1}
	if ok, _ = DifferByOneBit(a, d); ok {
		t.Error("DifferByOneBit false positive on different masks")
	}
}

func TestImplicantCovers(t *testing.T) {
	imp := Implicant{Value: 0b10, Mask: 0b01} // A is fixed to 1, B is don't care
	if !imp.Covers(0b10) || !imp.Covers(0b11) {
		t.Error("Covers should include 10 and 11")
	}
	if imp.Covers(0b00) || imp.Covers(0b01) {
		t.Error("Covers should not include 00 or 01")
	}
}

func TestGenerateSDNF(t *testing.T) {
	res := GenerateSDNF(2, []int{1, 3}, []string{"A", "B"})
	if res != "(!A & B) | (A & B)" {
		t.Errorf("Expected (!A & B) | (A & B), got %s", res)
	}
	if GenerateSDNF(2, []int{}, []string{"A", "B"}) != "0" {
		t.Error("Empty minterms should return 0")
	}
}

func TestMinimize(t *testing.T) {
	vars := []string{"A", "B", "C", "D"}
	// 1. All minterms -> "1"
	if Minimize(2, []int{0, 1, 2, 3}, nil, []string{"A", "B"}) != "1" {
		t.Error("All minterms should return 1")
	}
	// 2. No minterms -> "0"
	if Minimize(2, []int{}, nil, []string{"A", "B"}) != "0" {
		t.Error("No minterms should return 0")
	}
	// 3. Complex case
	res := Minimize(3, []int{0, 1, 2, 5, 6, 7}, []int{3}, []string{"A", "B", "C"})
	if len(res) == 0 {
		t.Error("Minimize failed on complex case")
	}

	// 4. Case to exercise findBestPrime
	res2 := Minimize(4, []int{0, 1, 2, 5, 6, 7, 8, 9, 10, 14}, nil, vars)
	if len(res2) == 0 {
		t.Error("Minimize failed on 4 var case")
	}
}

func TestMinimizeSimple(t *testing.T) {
	res := Minimize(2, []int{1, 3}, nil, []string{"A", "B"})
	if res != "(B)" {
		t.Errorf("Expected (B), got %s", res)
	}
}

func TestMinimizeWithDontCares(t *testing.T) {
	res := Minimize(2, []int{1}, []int{0}, []string{"A", "B"})
	if res != "(!A)" {
		t.Errorf("Expected (!A), got %s", res)
	}
}

func TestGetCovers(t *testing.T) {
	primes := []Implicant{{Value: 0, Mask: 1}, {Value: 2, Mask: 1}}
	covers := getCovers(primes, 0)
	if len(covers) != 1 || covers[0].Value != 0 {
		t.Error("getCovers failed")
	}
}

func TestAppendUnique(t *testing.T) {
	list := []Implicant{{Value: 1, Mask: 0}}
	list = appendUnique(list, Implicant{Value: 1, Mask: 0})
	if len(list) != 1 {
		t.Error("appendUnique failed for duplicate")
	}
	list = appendUnique(list, Implicant{Value: 2, Mask: 0})
	if len(list) != 2 {
		t.Error("appendUnique failed for new item")
	}
}

func TestFindBestPrime(t *testing.T) {
	primes := []Implicant{{Value: 0, Mask: 1}, {Value: 0, Mask: 3}}
	best := findBestPrime(primes, []int{0, 1, 2})
	if best.Mask != 3 {
		t.Error("findBestPrime failed to pick best")
	}
}

func TestFindPrimeImplicants(t *testing.T) {
	primes := findPrimeImplicants([]int{0, 1}, nil)
	expected := map[Implicant]struct{}{
		{Value: 0, Mask: 1}: {},
	}
	if len(primes) != len(expected) {
		t.Fatalf("Expected %d prime implicants, got %d", len(expected), len(primes))
	}
	for _, p := range primes {
		if _, ok := expected[p]; !ok {
			t.Fatalf("Unexpected prime implicant %+v", p)
		}
	}
}

func TestFindEssentialPrimes(t *testing.T) {
	primes := []Implicant{
		{Value: 0, Mask: 1}, // covers 0,1
		{Value: 2, Mask: 1}, // covers 2,3
		{Value: 0, Mask: 0}, // covers 0
	}
	essentials, remaining := findEssentialPrimes(primes, []int{0, 2})
	if len(essentials) != 1 || essentials[0].Value != 2 || essentials[0].Mask != 1 {
		t.Fatalf("Expected essential implicant {Value:2 Mask:1}, got %+v", essentials)
	}
	if len(remaining) != 1 || remaining[0] != 0 {
		t.Fatalf("Expected remaining [0], got %v", remaining)
	}
}

func TestCoverRemaining(t *testing.T) {
	primes := []Implicant{
		{Value: 0, Mask: 1}, // covers 0,1
		{Value: 2, Mask: 1}, // covers 2,3
	}
	remaining := []int{0, 2}
	solution := coverRemaining(remaining, primes, nil)
	if !coversAll(solution, remaining) {
		t.Fatalf("coverRemaining did not cover all minterms, solution=%v", solution)
	}
}

func TestFormat(t *testing.T) {
	vars := []string{"A", "B"}
	if formatImplicant(Implicant{Value: 0, Mask: 0}, 2, vars) != "(!A & !B)" {
		t.Error("formatImplicant failed")
	}
	if formatSolution([]Implicant{{Value: 0, Mask: 0}, {Value: 1, Mask: 0}}, 2, vars) != "(!A & !B) | (!A & B)" {
		t.Error("formatSolution failed")
	}
}

func coversAll(solution []Implicant, minterms []int) bool {
	for _, m := range minterms {
		covered := false
		for _, imp := range solution {
			if imp.Covers(m) {
				covered = true
				break
			}
		}
		if !covered {
			return false
		}
	}
	return true
}