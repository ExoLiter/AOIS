package evaluator

import "logical_calculator/internal/config"

// EvaluateRPN вычисляет значение функции на конкретном наборе
func EvaluateRPN(rpn []string, values map[string]bool) bool {
	var stack []bool
	for _, token := range rpn {
		stack = processRPNToken(token, stack, values)
	}
	if len(stack) == 0 {
		return false
	}
	return stack[0]
}

func processRPNToken(token string, stack []bool, values map[string]bool) []bool {
	if val, isVar := values[token]; isVar {
		return append(stack, val)
	}
	return applyOperator(token, stack)
}

func applyOperator(op string, stack []bool) []bool {
	if op == config.OpNot {
		return applyNot(stack)
	}
	return applyBinary(op, stack)
}

func applyNot(stack []bool) []bool {
	val := stack[len(stack)-1]
	return append(stack[:len(stack)-1], !val)
}

func applyBinary(op string, stack []bool) []bool {
	b := stack[len(stack)-1]
	a := stack[len(stack)-2]
	stack = stack[:len(stack)-2]

	var res bool
	switch op {
	case config.OpAnd:
		res = a && b
	case config.OpOr:
		res = a || b
	case config.OpImp:
		res = !a || b
	case config.OpEq:
		res = a == b
	}
	return append(stack, res)
}
