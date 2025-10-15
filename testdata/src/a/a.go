package a

import "testing"

// TestWithDefer demonstrates problematic use of defer in test functions.
// When t.Fatal or t.FailNow is called, they exit immediately via runtime.Goexit(),
// preventing deferred functions from running.
func TestWithDefer(t *testing.T) {
	defer cleanup() // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"

	// If t.Fatal is called here, the deferred cleanup() will not run
	// because t.Fatal calls runtime.Goexit() immediately
	if someCondition() {
		t.Fatal("test failed") // deferred cleanup() won't execute!
	}
}

// TestMultipleDefers shows multiple defer statements in a test
func TestMultipleDefers(t *testing.T) {
	defer func() { // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"
		// cleanup logic
	}()

	defer cleanup() // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"

	if someCondition() {
		t.FailNow() // none of the deferred functions will run
	}
}

// TestWithCleanup demonstrates the recommended pattern
func TestWithCleanup(t *testing.T) {
	t.Cleanup(cleanup) // No warning - this is the correct approach

	if someCondition() {
		t.Fatal("test failed") // cleanup will still run
	}
}

// TestMultipleCleanups shows multiple cleanup registrations
func TestMultipleCleanups(t *testing.T) {
	t.Cleanup(func() {
		// cleanup logic 1
	})

	t.Cleanup(cleanup) // No warning - this is the correct approach

	if someCondition() {
		t.Fatal("test failed") // all cleanups will still run
	}
}

// TestSubtestWithDefer shows defer is still problematic in subtests
func TestSubtestWithDefer(t *testing.T) {
	t.Run("subtest", func(t *testing.T) {
		defer cleanup() // want "use t.Cleanup\\(\\) instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow"

		if someCondition() {
			t.Fatal("failed")
		}
	})
}

// Helper functions should not trigger warnings
func helperWithDefer() {
	defer cleanup() // No warning - not a test function
}

func TestUsingHelper(t *testing.T) {
	helperWithDefer() // Should not trigger warning

	if someCondition() {
		t.Fatal("failed")
	}
}

func cleanup() {
	// cleanup logic
}

func someCondition() bool {
	return false
}
