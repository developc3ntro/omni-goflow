package main

// generate full docs with:
//
// go install github.com/developc3ntro/omni-goflow/cmd/docgen; docgen

import (
	"fmt"
	"os"

	"github.com/developc3ntro/omni-goflow/cmd/docgen/docs"
)

const (
	outputDir = "docs"
	localeDir = "locale"
)

func main() {
	if err := docs.Generate(".", outputDir, localeDir); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
