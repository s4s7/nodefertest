package a

import (
	"fmt"
	"testing"
)

// TestVariablesWithDeferName shows that variables containing "defer" are not flagged
func TestVariablesWithDeferName(t *testing.T) {
	// Variable names containing "defer"
	defer2 := "some value" // No warning - just a variable name
	deferCount := 0        // No warning - just a variable name
	shouldDefer := true    // No warning - just a variable name

	// String literals containing "defer"
	msg := "defer is a keyword"      // No warning - string literal
	comment := "use defer carefully" // No warning - string literal

	fmt.Println(defer2, deferCount, shouldDefer, msg, comment)

	// Actual defer statement - this SHOULD be flagged
	defer cleanup() // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"
}

// TestStructFieldsWithDefer shows struct fields named defer
func TestStructFieldsWithDefer(t *testing.T) {
	type Config struct {
		DeferTimeout int    // No warning - struct field
		UsesDefer    bool   // No warning - struct field
		deferMsg     string // No warning - struct field
	}

	cfg := Config{
		DeferTimeout: 100,
		UsesDefer:    false,
		deferMsg:     "defer related message",
	}

	fmt.Printf("%+v\n", cfg)

	// Actual defer - this SHOULD be flagged
	defer cleanup() // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"
}

// TestDeferInIfBlock shows defer inside if block
func TestDeferInIfBlock(t *testing.T) {
	if someCondition() {
		defer cleanup() // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"
	}
	t.Log("test")
}

// TestDeferInForLoop shows defer inside for loop
func TestDeferInForLoop(t *testing.T) {
	for i := 0; i < 3; i++ {
		defer cleanup() // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"
	}
	t.Log("test")
}

// TestDeferInSwitchCase shows defer inside switch statement
func TestDeferInSwitchCase(t *testing.T) {
	switch {
	case someCondition():
		defer cleanup() // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"
	default:
		defer func() {}() // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"
	}
	t.Log("test")
}

// TestDeferInSelect shows defer inside select statement
func TestDeferInSelect(t *testing.T) {
	ch := make(chan bool)
	select {
	case <-ch:
		defer cleanup() // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"
	default:
		defer func() {}() // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"
	}
	t.Log("test")
}

// TestNestedFuncLitWithoutTestingT shows nested function literal without *testing.T
func TestNestedFuncLitWithoutTestingT(t *testing.T) {
	defer cleanup() // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"

	// This function literal doesn't have *testing.T parameter
	// so defer inside it should not be flagged
	fn := func() {
		defer cleanup() // No warning - not a test function context
	}
	fn()
}

// TestNestedFuncLitWithTestingT shows nested function literal with *testing.T
func TestNestedFuncLitWithTestingT(t *testing.T) {
	defer cleanup() // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"

	// This function literal has *testing.T parameter
	// so defer inside it should be flagged
	fn := func(t *testing.T) {
		defer cleanup() // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"
	}
	fn(t)
}

// TestDeepNestedFuncLit shows deeply nested function literals
func TestDeepNestedFuncLit(t *testing.T) {
	t.Run("outer", func(t *testing.T) {
		defer cleanup() // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"

		t.Run("inner", func(t *testing.T) {
			defer cleanup() // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"

			// Function without *testing.T
			fn := func() {
				defer cleanup() // No warning - not a test function context
			}
			fn()
		})
	})
}

// BenchmarkWithDefer tests benchmark functions with defer
func BenchmarkWithDefer(b *testing.B) {
	defer cleanup() // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"

	for i := 0; i < b.N; i++ {
		// benchmark code
	}
}

// BenchmarkWithCleanup shows correct usage in benchmark
func BenchmarkWithCleanup(b *testing.B) {
	b.Cleanup(cleanup) // No warning - correct approach

	for i := 0; i < b.N; i++ {
		// benchmark code
	}
}

// BenchmarkSubBenchmarkWithDefer shows defer in sub-benchmark
func BenchmarkSubBenchmarkWithDefer(b *testing.B) {
	b.Run("sub", func(b *testing.B) {
		defer cleanup() // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"

		for i := 0; i < b.N; i++ {
			// benchmark code
		}
	})
}

// TestDeferWithPanicRecover shows defer with panic/recover
func TestDeferWithPanicRecover(t *testing.T) {
	defer func() { // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"
		if r := recover(); r != nil {
			t.Errorf("recovered: %v", r)
		}
	}()

	// test code
}

// TestAnonymousFunctionCall shows immediate function call
func TestAnonymousFunctionCall(t *testing.T) {
	func() {
		defer cleanup() // No warning - not a test function context
	}()

	// Function with *testing.T passed in
	func(t *testing.T) {
		defer cleanup() // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"
	}(t)
}

// TestDeferWithMultipleStatements shows defer with complex function
func TestDeferWithMultipleStatements(t *testing.T) {
	defer func() { // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"
		cleanup()
		// multiple statements
		if someCondition() {
			cleanup()
		}
	}()
}

// TestWithoutTestingT is a function starting with "Test" but no *testing.T parameter
// This should not trigger warnings
func TestWithoutTestingT() {
	defer cleanup() // No warning - no *testing.T parameter
}

// TestWithWrongParameterType shows Test function with different parameter
func TestWithWrongParameterType(s string) {
	defer cleanup() // No warning - parameter is not *testing.T
}

// TestMixedDeferAndCleanup shows mixed usage
func TestMixedDeferAndCleanup(t *testing.T) {
	defer cleanup()    // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"
	t.Cleanup(cleanup) // No warning - correct approach
	defer func() {}()  // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"
}

// TestDeferOfNamedFunction shows defer of different types of functions
func TestDeferOfNamedFunction(t *testing.T) {
	defer cleanup()             // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"
	defer helperWithDefer()     // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"
	defer t.Log("deferred log") // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"
}

// TestEmptyDefer shows defer with empty function
func TestEmptyDefer(t *testing.T) {
	defer func() {}() // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"
}

// TestDeferInGoroutine shows defer inside goroutine
func TestDeferInGoroutine(t *testing.T) {
	go func() {
		defer cleanup() // No warning - inside goroutine, not direct test function context
	}()
}

// TestDeferWithGoroutineAndTestingT shows goroutine with *testing.T parameter
func TestDeferWithGoroutineAndTestingT(t *testing.T) {
	go func(t *testing.T) {
		defer cleanup() // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"
	}(t)
}

// NotATestFunction is not a test function (doesn't start with "Test")
func NotATestFunction(t *testing.T) {
	defer cleanup() // No warning - not a test function
}

// Test is too short to be considered a test function
func Test(t *testing.T) {
	defer cleanup() // This might be an edge case - function name is exactly "Test"
}

// TestA is the shortest valid test function name
func TestA(t *testing.T) {
	defer cleanup() // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"
}
