package main

import (
	"github.com/vladik2801/loglint/logcheck"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
	unitchecker.Main(logcheck.Analyzer)
}
