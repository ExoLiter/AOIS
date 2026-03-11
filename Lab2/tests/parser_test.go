package tests

import (
	"logical_calculator/internal/evaluator"
	"logical_calculator/internal/parser"
	"testing"
)

func TestFormatInput(t *testing.T) {
	input := "!(!a -> !b) v c ~ d"
	expected := "!(!a>!b)|c=d"
	if res := parser.FormatInput(input); res != expected {
		t.Errorf("Expected %s, got %s", expected, res)
	}
}

func TestRPNAndEvaluation(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		values   map[string]bool
		expected bool
	}{
		{"AND true", "a&b", map[string]bool{"a": true, "b": true}, true},
		{"AND false", "a&b", map[string]bool{"a": true, "b": false}, false},
		{"Complex expression", "!(!a>!b)|c", map[string]bool{"a": true, "b": false, "c": false}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := parser.FormatInput(tt.expr)
			tokens := parser.Tokenize(formatted)
			rpn := parser.InfixToRPN(tokens)
			result := evaluator.EvaluateRPN(rpn, tt.values)
			if result != tt.expected {
				t.Errorf("Expression %s: expected %v, got %v", tt.expr, tt.expected, result)
			}
		})
	}
}
