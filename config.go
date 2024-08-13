package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func parseConfig() CommentsExtractorConfig {
	flagSet := flag.NewFlagSet("main", flag.ExitOnError)

	inputPath := flagSet.String("i", "./", "Root input path (file or folder)")
	outputFile := flagSet.String("o", "", "Output file. Required")
	inputPrefix := flagSet.String(
		"p",
		"",
		"Prefix going right after the beginning of Go multiline comment plus whitespace. Required",
	)
	outputPrefix := flagSet.String("c", "", "Prefix for comments in the output file. Optional")
	header := flagSet.String("h", "", "Header of the output file. Optional")

	flagSet.Usage = func() {
		fmt.Println("Extracts multiline comments with given prefix into a single file. Usage:")
		fmt.Println()
		flagSet.PrintDefaults()
	}

	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		flagSet.Usage()
		os.Exit(1)
	}

	if inputPrefix == nil {
		log.Fatal("specify input prefix -p. Usage: --help")
	}

	if outputFile == nil {
		log.Fatal("specify input prefix -p. Usage: --help")
	}

	if outputPrefix == nil {
		log.Fatal("specify input prefix -p. Usage: --help")
	}

	return CommentsExtractorConfig{
		InputCommentPattern: *inputPrefix,
		InputPath:           *inputPath,
		OutputPath:          *outputFile,
		OutputCommentPrefix: *outputPrefix,
		OutputFileHeader:    *header,
	}
}
