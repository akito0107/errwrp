package errwrp

import (
	"fmt"
	"go/ast"
	"go/token"
)

type Result struct {
	Position token.Position
	Pos      token.Pos
	Fname    string
}

func Check(aset *ParsedAST, fset *token.FileSet) ([]*Result, error) {
	res := check(aset, fset)
	return res, nil
}

func check(aset *ParsedAST, fset *token.FileSet) []*Result {
	var res []*Result

	ast.Inspect(aset.AST, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ReturnStmt:
			// may be return directly via function call
			if len(x.Results) < len(aset.AST.Type.Results.List) {
				res = append(res, &Result{
					Fname:    aset.FileName,
					Position: fset.Position(x.Pos()),
					Pos:      x.Pos(),
				})
				return true
			}
			expr := x.Results[aset.ErrOrd]
			if isNilError(expr) {
				return true
			}
			if ok, pkg := containsErrPkg(PkgErrors, aset.UsedErrorLikePackaged); ok {
				if isUsingPkgErrors(expr, pkg) {
					return true
				}
			}
			if ok, pkg := containsErrPkg(XErrors, aset.UsedErrorLikePackaged); ok {
				if isUsingXerrors(expr, pkg) {
					return true
				}
			}
			if mayUseOriginalError(expr) {
				return true
			}
			if isUsingFmtError(expr) {
				return true
			}
			res = append(res, &Result{
				Fname:    aset.FileName,
				Position: fset.Position(x.Pos()),
				Pos:      x.Pos(),
			})

		default:
			return true
		}
		return true
	})

	return res
}

// case of return nil, no problem
func isNilError(expr ast.Expr) bool {
	return fmt.Sprintf("%s", expr) == "nil"
}

// case of return errors.X may be safe
func isUsingPkgErrors(expr ast.Expr, ip ImportPath) bool {
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
	slctX := fmt.Sprintf("%s", slctexpr.X)
	if ip.Name != "" && ip.Name != slctX {
		return false
	} else if ip.Name == "" && slctX != "errors" {
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

func isUsingXerrors(expr ast.Expr, ip ImportPath) bool {
	var xerrorsMethods = []string{
		"Errorf", "New",
	}
	callexpr, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}
	slctexpr, ok := callexpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	slctX := fmt.Sprintf("%s", slctexpr.X)
	if ip.Name != "" && ip.Name != slctX {
		return false
	} else if ip.Name == "" && slctX != "xerrors" {
		return false
	}
	mname := fmt.Sprintf("%s", slctexpr.Sel.Name)
	for _, m := range xerrorsMethods {
		if mname == m {
			return true
		}
	}

	return false
}

func containsErrPkg(pkg ErrorPkg, arr []ImportPath) (bool, ImportPath) {
	for _, a := range arr {
		if pkg == a.Pkg {
			return true, a
		}
	}

	return false, ImportPath{}
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
