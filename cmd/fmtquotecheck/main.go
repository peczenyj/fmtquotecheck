package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/peczenyj/fmtquotecheck/analyzer"
)

func main() { singlechecker.Main(analyzer.New()) }
