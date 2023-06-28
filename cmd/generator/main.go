package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type CodeGenerator struct {
	Pattern string
}

func (g *CodeGenerator) GenerateSchemaFile(path string, outputPath string) error {
	// Get the absolute path of the specified directory
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	// Get a list of Go files in the specified directory and its subdirectories
	goFiles, err := g.getGoFiles(absPath)
	if err != nil {
		return err
	}

	// Open the output schema file in write mode
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Iterate over each Go file and extract the schema content
	for _, file := range goFiles {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(strings.NewReader(string(content)))
		scanner.Split(bufio.ScanLines)

		var lines []string
		var extracting bool

		for scanner.Scan() {
			line := scanner.Text()

			if extracting {
				if strings.HasPrefix(line, "*/") {
					extracting = false
					break
				}

				lines = append(lines, line)
			} else if strings.HasPrefix(line, fmt.Sprintf("/* %s", g.Pattern)) {
				extracting = true
				lines = append(lines, strings.TrimPrefix(line, fmt.Sprintf("/* %s", g.Pattern)))
			}
		}

		if err := scanner.Err(); err != nil {
			continue
		}

		if len(lines) == 0 {
			continue
		}

		// TODO: print relative file path of the input file

		schemaContent := strings.Join(lines, "\n")

		// Write the extracted schema content to the output file
		if _, err := outputFile.WriteString(schemaContent); err != nil {
			return fmt.Errorf("write line to output file: %w", err)
		}

		if _, err := outputFile.WriteString("\n"); err != nil {
			return fmt.Errorf("write line to output file: %w", err)
		}
	}

	return nil
}

// Helper function to get a list of Go files in the specified directory and its subdirectories
func (g *CodeGenerator) getGoFiles(dir string) ([]string, error) {
	var goFiles []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walking to %s: %w", path, err)
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			goFiles = append(goFiles, path)
		}

		return nil
	})

	return goFiles, err
}

func main() {
	g := CodeGenerator{
		Pattern: "SQL",
	}

	path := os.Args[1]
	outputPath := os.Args[2]

	g.GenerateSchemaFile(path, outputPath)
}
