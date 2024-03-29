package osexitcheck

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// Analyzer is a variable of the analysis type.The analyzer checks for a call to os.Exit()
var Analyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "check whether to use a direct call to os.Exit",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {

	for _, f := range pass.Files {

		ast.Inspect(f, func(node ast.Node) bool {
			// логика обработки
			callExpr, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}

			selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			ident, ok := selectorExpr.X.(*ast.Ident)
			if !ok {
				return true
			}

			if selectorExpr.Sel.Name == "Exit" && types.ExprString(selectorExpr) == "os.Exit" {
				pass.Reportf(ident.NamePos, "direct use of os.Exit found")
			}

			return true
		})
	}
	return nil, nil
}
