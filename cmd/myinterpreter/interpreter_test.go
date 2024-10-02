package main

import (
	"testing"
)

func TestInterpreter_Eval(t *testing.T) {
	code := "true"
	interpreter := NewInterpreter(code)

	value, err, _ := interpreter.Eval()
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
	interpreter := NewInterpreter(code)

	err, _ := interpreter.Run()
	if err != nil {
		t.Fatalf("interpreter.Run() error = %v", err)
	}
}

func TestInterpreter_Run_Print(t *testing.T) {
	code := `
		print "Hallo Welt!";
		print true;`
	interpreter := NewInterpreter(code)

	err, _ := interpreter.Run()
	if err != nil {
		t.Fatalf("interpreter.Run() error = %v", err)
	}
}

func TestInterpreter_Run_Fail(t *testing.T) {
	code := `
		print "the expression below is invalid";
		63 + "baz";
		print "this should not be printed";`
	interpreter := NewInterpreter(code)

	err, _ := interpreter.Run()
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
		}
		print quz;`
	interpreter := NewInterpreter(code)

	err, _ := interpreter.Run()
	if err != nil {
		t.Fatalf("interpreter.Run() error = %v", err)
	}
}
