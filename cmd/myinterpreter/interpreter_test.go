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

func TestInterpreter_WhileWithReturn(t *testing.T) {
	code := `
		fun f() {
			while (!false) return "ok";
		}
		print f();`

	interpreter := NewInterpreter(nil)

	err, _ := interpreter.Run(code)
	if err != nil {
		t.Fatalf("interpreter.Run() error = %v", err)
	}
}

func TestInterpreter_VarDecl(t *testing.T) {
	code := `
		var a = "outer";
		{
			var a = a;
		}`

	interpreter := NewInterpreter(nil)

	err, _ := interpreter.Run(code)
	if err == nil {
		t.Fatalf("expected interpreter error did not occur")
	}
}

func TestInterpreter_MutualRecursion(t *testing.T) {
	code := `
{
	   var threshold = 50;
 
   fun isEven(n) {
     if (n == 0) return true;
     if (n > threshold) return false;
     return isOdd(n - 1);
   }
 
   fun isOdd(n) {
     if (n == 0) return false;
     if (n > threshold) return false;
     return isEven(n - 1);
   }
 
   print isEven(5);
}`

	interpreter := NewInterpreter(nil)

	err, _ := interpreter.Run(code)
	if err != nil {
		t.Fatalf("interpreter.Run() error = %v", err)
	}
}

func TestInterpreter_VarDefWithExistingParameter(t *testing.T) {
	code := `
		fun foo(a) {
			var a = "hello";
		}`

	interpreter := NewInterpreter(nil)

	err, _ := interpreter.Run(code)
	if err == nil {
		t.Fatalf("expected interpreter error did not occur")
	}
}

func TestInterpreter_VarDeclWithSelf(t *testing.T) {
	code := `
		var a = "hello";
		var a = a;
		print a;
`
	interpreter := NewInterpreter(nil)

	err, _ := interpreter.Run(code)
	if err != nil {
		t.Fatalf("interpreter.Run() error = %v", err)
	}
}

func TestInterpreter_OuterVar(t *testing.T) {
	code := `
		fun makeCounter() {
			var i = 0;
			fun count() {
				i = i + 1;
				print i;
			}
			return count;
		}`

	interpreter := NewInterpreter(nil)

	err, _ := interpreter.Run(code)
	if err != nil {
		t.Fatalf("interpreter.Run() error = %v", err)
	}
}

func TestInterpreter_CallMethod(t *testing.T) {
	code := `
		class Foo {
			bar() {
				print "wunderbar!";
			}
		}
		var foo = Foo();
		foo.bar();`

	interpreter := NewInterpreter(nil)
	err, _ := interpreter.Run(code)
	if err != nil {
		t.Fatalf("interpreter.Run() error = %v", err)
	}
}

func TestInterpreter_SetProperty(t *testing.T) {
	code := `
		class Foo {
			bar(self) {
				print self.comment;
			}
		}
		var foo = Foo();
		foo.comment = "wunderbar!";
		foo.bar(foo);`

	interpreter := NewInterpreter(nil)
	err, _ := interpreter.Run(code)
	if err != nil {
		t.Fatalf("interpreter.Run() error = %v", err)
	}
}

func TestInterpreter_HigherOrderFunctions(t *testing.T) {
	code := `
		class Wizard {
			getSpellCaster() {
				fun castSpell() {
        			print this;
        			print "Casting spell as " + this.name;
     			}
 
    			// Functions are first-class objects in Lox
     			return castSpell;
   			}
 		}
 
		var wizard = Wizard();
		wizard.name = "Merlin";
 
		// Calling an instance method that returns a<|SPACE|>// function should work
		wizard.getSpellCaster()();`

	interpreter := NewInterpreter(nil)
	err, _ := interpreter.Run(code)
	if err != nil {
		t.Fatalf("interpreter.Run() error = %v", err)
	}
}
