package main

import (
	"fmt"
	"os"
)

func main() {
	g := CodeGenerator{
		Pattern:    inputPrefix,
		Path:       inputPath,
		outputPath: outputFile,
	}

	err := g.GenerateSchemaFile()
	if err != nil {
		fmt.Printf("Schema generation failed: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Schema generation completed!")
}
