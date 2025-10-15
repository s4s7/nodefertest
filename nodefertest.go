package nodefertest

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const doc = "nodefertest checks for the use of 'defer' in test functions, which can lead to unexpected behavior when functions like t.Fatal or t.FailNow are called, as they stop execution immediately and prevent deferred cleanup from running."

var Analyzer = &analysis.Analyzer{
	Name: "nodefertest",
	Doc:  doc,
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	// Iterate over all files
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			funcDecl, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}

			// Check if this is a test function
			if !isTestFunction(funcDecl) || !hasTestingTParam(funcDecl) {
				return true
			}

			// Check defer statements in this test function
			checkDeferInTestFunc(pass, funcDecl.Body)
			return false // Don't traverse into the function body again
		})
	}

	return nil, nil
}

// checkDeferInTestFunc recursively checks for defer statements in test functions
func checkDeferInTestFunc(pass *analysis.Pass, body *ast.BlockStmt) {
	ast.Inspect(body, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.DeferStmt:
			pass.Reportf(node.Defer,
				"use t.Cleanup() instead of defer in test functions to ensure cleanup runs even after t.Fatal/t.FailNow")
			return true
		case *ast.FuncLit:
			// Check if this function literal has a *testing.T parameter
			if hasFuncLitTestingTParam(node) {
				// Recursively check this function literal
				checkDeferInTestFunc(pass, node.Body)
			}
			// Don't traverse into this function literal from here
			// (we already handled it above if it has *testing.T param)
			return false
		}
		return true
	})
}

// hasFuncLitTestingTParam checks if the function literal has a *testing.T parameter
func hasFuncLitTestingTParam(funcLit *ast.FuncLit) bool {
	if funcLit.Type == nil || funcLit.Type.Params == nil {
		return false
	}

	for _, field := range funcLit.Type.Params.List {
		starExpr, ok := field.Type.(*ast.StarExpr)
		if !ok {
			continue
		}

		selectorExpr, ok := starExpr.X.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		ident, ok := selectorExpr.X.(*ast.Ident)
		if !ok {
			continue
		}

		// Check if it's testing.T or testing.B
		if ident.Name == "testing" && (selectorExpr.Sel.Name == "T" || selectorExpr.Sel.Name == "B") {
			return true
		}
	}

	return false
}

// isTestFunction checks if the function is a test function
func isTestFunction(funcDecl *ast.FuncDecl) bool {
	name := funcDecl.Name.Name
	// Test functions start with "Test", "Benchmark", or "Example"
	if len(name) > 4 && name[:4] == "Test" {
		return true
	}
	if len(name) > 9 && name[:9] == "Benchmark" {
		return true
	}
	return false
}

// hasTestingTParam checks if the function has a *testing.T parameter
func hasTestingTParam(funcDecl *ast.FuncDecl) bool {
	if funcDecl.Type.Params == nil || len(funcDecl.Type.Params.List) == 0 {
		return false
	}

	for _, field := range funcDecl.Type.Params.List {
		starExpr, ok := field.Type.(*ast.StarExpr)
		if !ok {
			continue
		}

		selectorExpr, ok := starExpr.X.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		ident, ok := selectorExpr.X.(*ast.Ident)
		if !ok {
			continue
		}

		// Check if it's testing.T or testing.B
		if ident.Name == "testing" && (selectorExpr.Sel.Name == "T" || selectorExpr.Sel.Name == "B") {
			return true
		}
	}

	return false
}
