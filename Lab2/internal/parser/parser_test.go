package parser

import (
	"reflect"
	"testing"
)

func TestFormatInputVariants(t *testing.T) {
	input := " a v b ∨ c -> d → e ~ f "
	got := FormatInput(input)
	want := "a|b|c>d>e=f"
	if got != want {
		t.Fatalf("FormatInput mismatch: got %q want %q", got, want)
	}
}

func TestTokenize(t *testing.T) {
	input := "a&b"
	want := []string{"a", "&", "b"}
	if got := Tokenize(input); !reflect.DeepEqual(got, want) {
		t.Fatalf("Tokenize mismatch: got %#v want %#v", got, want)
	}
}

func TestInfixToRPNPrecedence(t *testing.T) {
	tokens := []string{"a", "|", "b", "&", "c"}
	got := InfixToRPN(tokens)
	want := []string{"a", "b", "c", "&", "|"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("InfixToRPN mismatch: got %#v want %#v", got, want)
	}
}

func TestInfixToRPNParenthesesAndNot(t *testing.T) {
	tokens := []string{"!", "(", "a", "&", "b", ")"}
	got := InfixToRPN(tokens)
	want := []string{"a", "b", "&", "!"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("InfixToRPN mismatch: got %#v want %#v", got, want)
	}
}

func TestValidateExpression(t *testing.T) {
	tests := []struct {
		name   string
		tokens []string
		valid  bool
	}{
		{"empty", []string{}, false},
		{"extra closing", []string{")"}, false},
		{"missing closing", []string{"(", "a"}, false},
		{"unknown token", []string{"a", "$"}, false},
		{"bad sequence", []string{"a", "("}, false},
		{"binary before )", []string{"a", "&", ")"}, false},
		{"valid", []string{"(", "a", "&", "b", ")", "|", "!", "c"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateExpression(tt.tokens)
			if tt.valid && err != nil {
				t.Fatalf("expected valid, got error: %v", err)
			}
			if !tt.valid && err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
	}
}

func TestGetPriorityDefaultAndIsVariable(t *testing.T) {
	if p := getPriority("?"); p != 0 {
		t.Fatalf("expected default priority 0, got %d", p)
	}
	if isVariable("f") {
		t.Fatalf("expected 'f' to be non-variable")
	}
	if !isVariable("a") {
		t.Fatalf("expected 'a' to be variable")
	}
}
