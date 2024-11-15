package main

import (
	"testing"
)

func TestInterpreter_Eval(t *testing.T) {
	code := "true"
	interpreter := NewInterpreter(nil)

	value, err, _ := interpreter.Eval(code)
	if err != nil {
		t.Fatalf("interpreter.Eval() error = %v", err)
	}

	expectedValue := NewBooleanValue(true)
	if value.getType() != expectedValue.getType() {
		t.Fatalf("Expected %s, got %s", expectedValue, value)
	}

}

func TestInterpreter_Run_VarDecl(t *testing.T) {
	code := `
		var a = "foo";
		print a;`
	interpreter := NewInterpreter(nil)

	err, _ := interpreter.Run(code)
	if err != nil {
		t.Fatalf("interpreter.Run() error = %v", err)
	}
}

func TestInterpreter_Run_Print(t *testing.T) {
	code := `
		print "Hallo Welt!";
		print true;`
	interpreter := NewInterpreter(nil)

	err, _ := interpreter.Run(code)
	if err != nil {
		t.Fatalf("interpreter.Run() error = %v", err)
	}
}

func TestInterpreter_Run_Fail(t *testing.T) {
	code := `
		print "the expression below is invalid";
		63 + "baz";
		print "this should not be printed";`
	interpreter := NewInterpreter(nil)

	err, _ := interpreter.Run(code)
	if err == nil {
		t.Fatalf("expected interpreter error did not occur")
	}
}

func TestInterpreter_Run_Block(t *testing.T) {
	code := `
		{
			var bar = "outer bar";
			var quz = "outer quz";
			{
				bar = "modified bar";
				var quz = "inner quz";
				print bar;
				print quz;
			}
			print bar;
			print quz;
		}`
	interpreter := NewInterpreter(nil)

	err, _ := interpreter.Run(code)
	if err != nil {
		t.Fatalf("interpreter.Run() error = %v", err)
	}
}

func TestInterpreter_Call_Builtin(t *testing.T) {
	code := `print clock();`
	interpreter := NewInterpreter(nil)

	err, _ := interpreter.Run(code)
	if err != nil {
		t.Fatalf("interpreter.Run() error = %v", err)
	}
}

func TestInterpreter_FunctionDef(t *testing.T) {
	code := `
	fun f3(a, b, c) { print a + b + c; }
	f3(27, 27, 27);
	`
	interpreter := NewInterpreter(nil)

	err, _ := interpreter.Run(code)
	if err != nil {
		t.Fatalf("interpreter.Run() error = %v", err)
	}
}
