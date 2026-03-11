package parser

import (
	"logical_calculator/internal/config"
	"strings"
)

// FormatInput очищает ввод от пробелов и нормализует операторы
func FormatInput(input string) string {
	input = strings.ReplaceAll(input, " ", "")
	input = strings.ReplaceAll(input, "v", config.OpOr)
	input = strings.ReplaceAll(input, "∨", config.OpOr)
	input = strings.ReplaceAll(input, "->", config.OpImp)
	input = strings.ReplaceAll(input, "→", config.OpImp)
	input = strings.ReplaceAll(input, "~", config.OpEq)
	return input
}

// Tokenize разбивает выражение на токены
func Tokenize(input string) []string {
	var tokens []string
	for _, char := range input {
		tokens = append(tokens, string(char))
	}
	return tokens
}

// InfixToRPN переводит инфиксную запись в обратную польскую нотацию
func InfixToRPN(tokens []string) []string {
	var rpn, stack []string
	for _, token := range tokens {
		rpn, stack = processTokenForRPN(token, rpn, stack)
	}
	for len(stack) > 0 {
		rpn = append(rpn, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	return rpn
}

func processTokenForRPN(token string, rpn, stack []string) ([]string, []string) {
	if isVariable(token) {
		return append(rpn, token), stack
	}
	if token == config.OpLPr {
		return rpn, append(stack, token)
	}
	if token == config.OpRPr {
		return flushUntilLeftParen(rpn, stack)
	}
	return handleOperatorForRPN(token, rpn, stack)
}

func flushUntilLeftParen(rpn, stack []string) ([]string, []string) {
	for len(stack) > 0 && stack[len(stack)-1] != config.OpLPr {
		rpn = append(rpn, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	if len(stack) > 0 {
		stack = stack[:len(stack)-1] // Удаляем "("
	}
	return rpn, stack
}

func handleOperatorForRPN(op string, rpn, stack []string) ([]string, []string) {
	for len(stack) > 0 && getPriority(stack[len(stack)-1]) >= getPriority(op) {
		rpn = append(rpn, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	return rpn, append(stack, op)
}

func getPriority(op string) int {
	switch op {
	case config.OpNot:
		return config.PriorityNot
	case config.OpAnd:
		return config.PriorityAnd
	case config.OpOr:
		return config.PriorityOr
	case config.OpImp:
		return config.PriorityImp
	case config.OpEq:
		return config.PriorityEq
	}
	return config.PriorityDef
}

func isVariable(token string) bool {
	return token >= "a" && token <= "e"
}
