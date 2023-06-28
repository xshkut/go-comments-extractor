package main

import (
	"fmt"
	"os"
)

func main() {
	g := commentsExtractor{
		inputPattern:        inputPrefix,
		inputPath:           inputPath,
		outputPath:          outputFile,
		outputCommentPrefix: outputPrefix,
		header:              header,
	}

	err := g.ExtractComments()
	if err != nil {
		fmt.Printf("Schema generation failed: %s\n", err)
		os.Exit(1)
	}
}
