package minimization

import (
	"fmt"
	"logical_calculator/internal/models"
	"strings"
)

func MinimizeTabularCalcSDNF(table []models.Row, vars []string) string {
	minterms := BuildMinterms(table, vars)
	return runTabularCalcMethod(minterms, vars, FormatTermSDNF, FormatTermsSum, "v")
}

func MinimizeTabularCalcSKNF(table []models.Row, vars []string) string {
	maxterms := BuildMaxterms(table, vars)
	return runTabularCalcMethod(maxterms, vars, FormatTermSKNF, FormatTermsProd, "&")
}

func runTabularCalcMethod(terms []models.Term, vars []string,
	formatTerm func(models.Term, []string) string,
	formatFinal func([]models.Term, []string) string, op string) string {

	if len(terms) == 0 {
		return "Нет импликант для таблицы."
	}

	var sb strings.Builder
	sb.WriteString("Этап склеивания (аналогично расчетному методу):\n")
	primes := PerformGluingSteps(&sb, terms, vars, formatTerm, op)

	sb.WriteString("\nПостроение таблицы (Расчетно-табличный):\n")
	sb.WriteString(buildQuineTable(primes, terms, vars, formatTerm))

	final := RemoveRedundantCalc(nil, primes, terms, vars, formatTerm) // nil чтобы не дублировать лог
	sb.WriteString("\nУбираем лишние импликанты и получаем:\n")
	sb.WriteString(formatFinal(final, vars) + "\n")

	return sb.String()
}

func buildQuineTable(primes, baseTerms []models.Term, vars []string, formatTerm func(models.Term, []string) string) string {
	var sb strings.Builder

	// 1. Динамически вычисляем максимальную ширину для первого столбца
	maxHeaderLen := len("Импликанты")
	for _, p := range primes {
		termStr := formatTerm(p, vars)
		if len(termStr) > maxHeaderLen {
			maxHeaderLen = len(termStr)
		}
	}
	maxHeaderLen += 2 // Добавляем отступ для красоты

	// 2. Устанавливаем узкую ширину для столбцов с индексами
	colWidth := 4

	// 3. Шапка таблицы (выводим ТОЛЬКО индексы минтермов/макстермов)
	sb.WriteString(fmt.Sprintf("%-*s ", maxHeaderLen, "Импликанты"))
	for _, m := range baseTerms {
		sb.WriteString(fmt.Sprintf("| %-*d ", colWidth, m.Value))
	}
	sb.WriteString("|\n")

	// 4. Линия-разделитель (считаем точную длину)
	lineLen := maxHeaderLen + 1 + len(baseTerms)*(colWidth+3)
	sb.WriteString(strings.Repeat("-", lineLen) + "\n")

	// 5. Тело таблицы
	for _, p := range primes {
		sb.WriteString(fmt.Sprintf("%-*s ", maxHeaderLen, formatTerm(p, vars)))
		for _, m := range baseTerms {
			if covers(p, m) {
				sb.WriteString(fmt.Sprintf("| %-*s ", colWidth, " X")) // Ставим крестик
			} else {
				sb.WriteString(fmt.Sprintf("| %-*s ", colWidth, "")) // Пустота
			}
		}
		sb.WriteString("|\n")
	}

	return sb.String()
}
