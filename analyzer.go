package errwrp

import (
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name:             "mustwrap",
	Doc:              "check for just returning error without errors.Wrap",
	RunDespiteErrors: false,
	Run:              run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	var decls []*ParsedAST
	for _, f := range pass.Files {
		decls = append(decls, parse("", f)...)
	}

	var res []*Result
	for _, d := range decls {
		res = append(res, check(d, pass.Fset)...)
	}

	for _, r := range res {
		pass.Reportf(r.Pos, "should be use errors.Wrap() or errors.Wrapf()")
	}

	return nil, nil
}
