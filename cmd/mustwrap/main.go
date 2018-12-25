package main

import (
	"github.com/akito0107/errwrp"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(errwrp.Analyzer)
}
