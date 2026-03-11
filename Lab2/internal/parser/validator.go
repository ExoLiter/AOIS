package parser

import (
	"errors"
	"fmt"
	"logical_calculator/internal/config"
)

// ValidateExpression проверяет массив токенов на синтаксические ошибки
func ValidateExpression(tokens []string) error {
	if len(tokens) == 0 {
		return errors.New("введено пустое выражение")
	}

	balance := 0
	for i, token := range tokens {
		// 1. Проверка баланса скобок
		if token == config.OpLPr {
			balance++
		} else if token == config.OpRPr {
			balance--
			if balance < 0 {
				return errors.New("ошибка синтаксиса: лишняя закрывающая скобка ')'")
			}
		}

		// 2. Проверка на неизвестные символы
		if !isValidToken(token) {
			return fmt.Errorf("неизвестный символ: '%s'", token)
		}

		// 3. Проверка последовательности символов
		if i > 0 {
			if err := checkSequence(tokens[i-1], token); err != nil {
				return err
			}
		}
	}

	if balance != 0 {
		return errors.New("ошибка синтаксиса: не хватает закрывающей скобки")
	}

	return nil
}

func isValidToken(token string) bool {
	return isVariable(token) || isBinaryOperator(token) ||
		token == config.OpNot || token == config.OpLPr || token == config.OpRPr
}

func isBinaryOperator(token string) bool {
	switch token {
	case config.OpAnd, config.OpOr, config.OpImp, config.OpEq:
		return true
	}
	return false
}

func checkSequence(prev, curr string) error {
	if (isBinaryOperator(prev) || prev == config.OpLPr) && (isBinaryOperator(curr) || curr == config.OpRPr) {
		return fmt.Errorf("ошибка синтаксиса: неверная последовательность '%s%s'", prev, curr)
	}
	if (isVariable(prev) || prev == config.OpRPr) && (isVariable(curr) || curr == config.OpLPr || curr == config.OpNot) {
		return fmt.Errorf("ошибка синтаксиса: пропущен оператор между '%s' и '%s'", prev, curr)
	}
	return nil
}
