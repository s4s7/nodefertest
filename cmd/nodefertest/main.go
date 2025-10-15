package main

import (
	"github.com/s4s7/nodefertest"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(nodefertest.Analyzer) }
