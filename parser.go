package errwrp

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"

	"github.com/google/logger"
	"github.com/pkg/errors"
)

func init() {
	logger.Init("errwrp", true, true, os.Stdout)
}

type ImportPath struct {
	Name string
	Pkg  ErrorPkg
}

type ParsedAST struct {
	ErrOrd                int
	AST                   *ast.FuncDecl
	FileName              string
	UsedErrorLikePackaged []ImportPath
}

func Parse(r io.Reader, fname string) ([]*ParsedAST, *token.FileSet, ast.CommentMap, error) {
	src, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "parser: ioutil/ReadAll")
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", string(src), parser.ParseComments)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "parser: parser/ParseFile")
	}

	decls := parse(fname, f)

	return decls, fset, ast.NewCommentMap(fset, f, f.Comments), nil
}

func parse(fname string, f *ast.File) []*ParsedAST {
	var decls []*ParsedAST
	var usedPackages []ImportPath

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ImportSpec:
			ep := FromPkgName(x.Path.Value)
			if ep == Unknown {
				return true
			}
			p := ImportPath{
				Pkg: ep,
			}
			if x.Name != nil {
				p.Name = x.Name.String()
			}
			usedPackages = append(usedPackages, p)
		case *ast.FuncDecl:
			if x.Type.Results == nil {
				return true
			}
			resList := x.Type.Results.List
			for i := 0; i < len(resList); i++ {
				if !containsErrorType(resList[i].Type) {
					continue
				}
				if isNamedReturn(resList[i]) {
					logger.Warningf("named return currently not supported: %v", x.Name)
					continue
				}
				p := &ParsedAST{
					ErrOrd:                i,
					AST:                   x,
					FileName:              fname,
					UsedErrorLikePackaged: usedPackages,
				}
				decls = append(decls, p)
			}
			return true
		default:
			return true
		}

		return true
	})
	return decls
}

func containsErrorType(e ast.Expr) bool {
	return fmt.Sprint(e) == "error"
}

func isNamedReturn(f *ast.Field) bool {
	return len(f.Names) != 0
}
