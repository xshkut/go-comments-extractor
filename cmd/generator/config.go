package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	inputPath    string
	outputFile   string
	inputPrefix  string
	outputPrefix string
	header       string
)

func init() {
	flag.StringVar(&inputPath, "i", "./", "Root input path (file or folder). Default: ./")
	flag.StringVar(&outputFile, "o", "", "Output file. Required")
	flag.StringVar(&inputPrefix, "p", "", "Prefix going right after the beginning of Go multiline comment plus whitespace. Required")
	flag.StringVar(&outputPrefix, "c", "", "Prefix for comments in the output file. Optional")
	flag.StringVar(&header, "h", "", "Header of the output file. Optional")

	flag.Usage = func() {
		fmt.Println("Extracts multiline comments with given prefix into a single file. Usage:")
		fmt.Println()
		flag.PrintDefaults()
	}

	flag.Parse()

	if inputPrefix == "" {
		log.Fatal("specify input prefix -p. Usage: --help")
	}

	if outputFile == "" {
		log.Fatal("specify input prefix -p. Usage: --help")
	}
}
