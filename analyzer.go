package errwrp

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name:             "mustwrap",
	Doc:              "check for just returning error without errors.Wrap",
	RunDespiteErrors: false,
	Run:              run,
}

type decl struct {
	Comment ast.CommentMap
	Decls   []*ParsedAST
}

func run(pass *analysis.Pass) (interface{}, error) {
	var decls []decl

	for _, f := range pass.Files {
		commMap := ast.NewCommentMap(pass.Fset, f, f.Comments)
		decls = append(decls,
			decl{
				Comment: commMap,
				Decls:   parse("", f),
			})
	}

	var res []*Result
	for _, d := range decls {
		for _, a := range d.Decls {
			res = append(res, check(a, pass.Fset, d.Comment)...)
		}
	}

	for _, r := range res {
		pass.Reportf(r.Pos, "should use errors.Wrap() or errors.Wrapf()")
	}

	return nil, nil
}
