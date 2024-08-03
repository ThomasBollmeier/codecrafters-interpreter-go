package main

import (
	"testing"
)

func TestInterpreter_Eval(t *testing.T) {
	code := "true"
	interpreter := NewInterpreter(code)

	value, err := interpreter.Eval()
	if err != nil {
		t.Fatalf("interpreter.Eval() error = %v", err)
	}

	expectedValue := NewBooleanValue(true)
	if value.getType() != expectedValue.getType() {
		t.Fatalf("Expected %s, got %s", expectedValue, value)
	}

}
