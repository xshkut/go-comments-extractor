package main

import (
	"fmt"
	"os"
)

func main() {
	cfg := parseConfig()

	err := NewCommentsExtractor(cfg).ExtractComments()
	if err != nil {
		fmt.Printf("Schema generation failed: %s\n", err)
		os.Exit(1)
	}
}
