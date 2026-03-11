package main

import (
	"bufio"
	"fmt"
	"logical_calculator/internal/derivatives"
	"logical_calculator/internal/forms"
	"logical_calculator/internal/minimization"
	"logical_calculator/internal/parser"
	"logical_calculator/internal/postclass"
	"logical_calculator/internal/truthtable"
	"logical_calculator/internal/zhegalkin"
	"os"
)

func main() {
	fmt.Println("Введите логическую функцию (например: !(!a->!b)v c):")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return
	}
	input := scanner.Text()

	// 1. Форматирование, токенизация и валидация
	formatted := parser.FormatInput(input)
	tokens := parser.Tokenize(formatted)

	if err := parser.ValidateExpression(tokens); err != nil {
		fmt.Printf("Ошибка во введенной формуле: %v\n", err)
		return
	}

	rpn := parser.InfixToRPN(tokens)

	// 2. Таблица истинности
	vars := truthtable.ExtractVariables(rpn)
	table := truthtable.GenerateTable(rpn, vars)

	fmt.Println("\n--- Таблица истинности ---")
	for _, v := range vars {
		fmt.Printf("%s\t", v)
	}
	fmt.Println("F")
	for _, row := range table {
		for _, v := range vars {
			if row.Values[v] {
				fmt.Print("1\t")
			} else {
				fmt.Print("0\t")
			}
		}
		if row.Result {
			fmt.Println("1")
		} else {
			fmt.Println("0")
		}
	}

	// 3, 4, 5. Формы
	sdnfNums, sknfNums := forms.NumericForms(table)
	fmt.Printf("\nСДНФ: %s\nЧисловая СДНФ: %v\n", forms.BuildSDNF(table, vars), sdnfNums)
	fmt.Printf("СКНФ: %s\nЧисловая СКНФ: %v\n", forms.BuildSKNF(table, vars), sknfNums)
	fmt.Printf("Индексная форма: %d\n", forms.IndexForm(table))

	// 6. Классы Поста
	fmt.Println("\n--- Классы Поста ---")
	fmt.Printf("T0: %v, T1: %v, S: %v, M: %v, L: %v\n",
		postclass.IsT0(table), postclass.IsT1(table),
		postclass.IsSelfDual(table), postclass.IsMonotonic(table), postclass.IsLinear(table))

	if postclass.IsFunctionallyComplete(table) {
		fmt.Println("Данная логическая функция является функционально полной")
	} else {
		fmt.Println("Данная логическая функция не является функционально полной")
	}

	// 7. Полином Жегалкина
	fmt.Printf("\nПолином Жегалкина: %s\n", zhegalkin.BuildPolynomial(table, vars))

	// 8. Фиктивные переменные
	dummies := derivatives.FindDummyVariables(table, vars)
	fmt.Printf("\nФиктивные переменные: %v\n", dummies)

	// 9. Булева дифференциация (Все частные и смешанные до 4 порядка)
	fmt.Println("\n--- Булева дифференциация ---")
	fmt.Println(derivatives.GenerateAllDerivatives(table, vars))

	// 10, 11, 12. МИНИМИЗАЦИЯ СДНФ
	fmt.Println("\n============= МИНИМИЗАЦИЯ СДНФ =============")

	fmt.Println("\n--- Расчетный метод (СДНФ) ---")
	fmt.Println(minimization.MinimizeCalculationSDNF(table, vars))

	fmt.Println("\n--- Расчетно-табличный метод (СДНФ) ---")
	fmt.Println(minimization.MinimizeTabularCalcSDNF(table, vars))

	fmt.Println("\n--- Табличный метод / Карта Карно (СДНФ) ---")
	fmt.Println(minimization.MinimizeKarnaughSDNF(table, vars))

	// 10, 11, 12. МИНИМИЗАЦИЯ СКНФ
	fmt.Println("\n============= МИНИМИЗАЦИЯ СКНФ =============")

	fmt.Println("\n--- Расчетный метод (СКНФ) ---")
	fmt.Println(minimization.MinimizeCalculationSKNF(table, vars))

	fmt.Println("\n--- Расчетно-табличный метод (СКНФ) ---")
	fmt.Println(minimization.MinimizeTabularCalcSKNF(table, vars))

	fmt.Println("\n--- Табличный метод / Карта Карно (СКНФ) ---")
	fmt.Println(minimization.MinimizeKarnaughSKNF(table, vars))
}
