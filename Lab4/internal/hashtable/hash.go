package hashtable

import (
	"strings"
	"unicode"
)

func computeV(key string) (int, error) {
	letters, err := firstTwoLetters(key)
	if err != nil {
		return 0, err
	}
	return computeVFromLetters(letters)
}

func computeH(v int, size int, base int) int {
	return (v % size) + base
}

func computeVFromLetters(letters []rune) (int, error) {
	first, second := letters[0], letters[1]
	if idx1, ok := russianIndex(first); ok {
		idx2, ok2 := russianIndex(second)
		if !ok2 {
			return 0, ErrKeyAlphabet
		}
		return idx1*RussianBase + idx2, nil
	}
	if idx1, ok := latinIndex(first); ok {
		idx2, ok2 := latinIndex(second)
		if !ok2 {
			return 0, ErrKeyAlphabet
		}
		return idx1*LatinBase + idx2, nil
	}
	return 0, ErrKeyAlphabet
}

func firstTwoLetters(key string) ([]rune, error) {
	trimmed := strings.TrimSpace(key)
	letters := make([]rune, 0, 2)
	for _, r := range trimmed {
		if !unicode.IsLetter(r) {
			continue
		}
		letters = append(letters, unicode.ToUpper(r))
		if len(letters) == 2 {
			break
		}
	}
	if len(letters) < 2 {
		return nil, ErrKeyInvalid
	}
	return letters, nil
}

func normalizeKey(key string) (string, error) {
	letters, err := firstTwoLetters(key)
	if err != nil {
		return "", err
	}
	upper := strings.ToUpper(strings.TrimSpace(key))
	_, err = computeVFromLetters(letters)
	if err != nil {
		return "", err
	}
	return upper, nil
}

func russianIndex(r rune) (int, bool) {
	return runeIndex(RussianAlphabet, r)
}

func latinIndex(r rune) (int, bool) {
	return runeIndex(LatinAlphabet, r)
}

func runeIndex(alphabet string, r rune) (int, bool) {
	index := 0
	for _, ar := range alphabet {
		if ar == r {
			return index, true
		}
		index++
	}
	return 0, false
}
