package errwrp

import (
	"fmt"
	"go/ast"
	"go/token"
)

type Result struct {
	Pos   token.Position
	Fname string
}

func Check(aset *ParsedAST, fset *token.FileSet) ([]*Result, error) {
	var res []*Result

	ast.Inspect(aset.AST, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ReturnStmt:
			// may be return directly via function call
			if len(x.Results) < len(aset.AST.Type.Results.List) {
				res = append(res, &Result{
					Fname: aset.FileName,
					Pos:   fset.Position(x.Pos()),
				})
				return true
			}
			expr := x.Results[aset.ErrOrd]
			if isNilError(expr) {
				return true
			}
			if isUsingPkgErrors(expr) {
				return true
			}
			if mayUseOriginalError(expr) {
				return true
			}
			if isUsingFmtError(expr) {
				return true
			}
			res = append(res, &Result{
				Fname: aset.FileName,
				Pos:   fset.Position(x.Pos()),
			})

		default:
			return true
		}
		return true
	})

	return res, nil
}

// case of return nil, no problem
func isNilError(expr ast.Expr) bool {
	return fmt.Sprintf("%s", expr) == "nil"
}

// case of return errors.X may be safe
func isUsingPkgErrors(expr ast.Expr) bool {
	var pkgErrorMethods = []string{
		"Errorf", "New", "WithMessage", "WithMessagef",
		"WithStack", "Wrap", "Wrapf",
	}
	callexpr, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}
	slctexpr, ok := callexpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if fmt.Sprintf("%s", slctexpr.X) != "errors" {
		return false
	}
	mname := fmt.Sprintf("%s", slctexpr.Sel.Name)
	for _, m := range pkgErrorMethods {
		if mname == m {
			return true
		}
	}

	return false
}

// case of custom errors such as &myerror{}, no problem
func mayUseOriginalError(expr ast.Expr) bool {
	_, ok := expr.(*ast.UnaryExpr)
	return ok
}

func isUsingFmtError(expr ast.Expr) bool {
	callexpr, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}
	slctexpr, ok := callexpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if fmt.Sprintf("%s", slctexpr.X) != "fmt" {
		return false
	}
	mname := fmt.Sprintf("%s", slctexpr.Sel.Name)
	return mname == "Errorf"

}
