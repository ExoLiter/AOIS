package evaluator

import "testing"

func TestEvaluateRPNBasicOps(t *testing.T) {
	rpn := []string{"a", "b", "&"}
	values := map[string]bool{"a": true, "b": false}
	if got := EvaluateRPN(rpn, values); got != false {
		t.Fatalf("expected false, got %v", got)
	}

	rpn = []string{"a", "!"}
	values = map[string]bool{"a": false}
	if got := EvaluateRPN(rpn, values); got != true {
		t.Fatalf("expected true, got %v", got)
	}
}

func TestEvaluateRPNAllBinaryOps(t *testing.T) {
	tests := []struct {
		rpn      []string
		values   map[string]bool
		expected bool
	}{
		{[]string{"a", "b", "|"}, map[string]bool{"a": false, "b": true}, true},
		{[]string{"a", "b", ">"}, map[string]bool{"a": true, "b": false}, false},
		{[]string{"a", "b", "="}, map[string]bool{"a": true, "b": true}, true},
	}
	for _, tt := range tests {
		if got := EvaluateRPN(tt.rpn, tt.values); got != tt.expected {
			t.Fatalf("unexpected result for %v: got %v want %v", tt.rpn, got, tt.expected)
		}
	}
}

func TestEvaluateRPNEmpty(t *testing.T) {
	if got := EvaluateRPN([]string{}, map[string]bool{}); got != false {
		t.Fatalf("expected false on empty RPN, got %v", got)
	}
}
