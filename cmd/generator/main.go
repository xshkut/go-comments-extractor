package main

import (
	"fmt"
	"os"
)

func main() {
	g := codeGenerator{
		inputPattern:        inputPrefix,
		inputPath:           inputPath,
		outputPath:          outputFile,
		outputCommentPrefix: outputPrefix,
		header:              header,
	}

	err := g.GenerateSchemaFile()
	if err != nil {
		fmt.Printf("Schema generation failed: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Schema generation completed!")
}
