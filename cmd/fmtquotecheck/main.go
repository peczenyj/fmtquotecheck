package main

import (
	"log"

	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/peczenyj/fmtquotecheck/analyzer"
)

func main() {
	instance, err := analyzer.New()
	if err != nil {
		log.Fatalf("unexpected error while building analyzer: %v", err)
	}

	singlechecker.Main(instance)
}
