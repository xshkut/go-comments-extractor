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
	Pattern    string
	Path       string
	outputPath string
}

func (g *CodeGenerator) GenerateSchemaFile() error {
	// Get the absolute path of the specified directory
	inputAbsPath, err := filepath.Abs(g.Path)
	if err != nil {
		return err
	}

	outputAbsPath, err := filepath.Abs(g.outputPath)
	if err != nil {
		return err
	}

	// Get a list of Go files in the specified directory and its subdirectories
	goFiles, err := g.getGoFiles(inputAbsPath)
	if err != nil {
		return err
	}

	// Open the output schema file in write mode
	outputFile, err := os.Create(g.outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Iterate over each Go file and extract the schema content
	for _, filePath := range goFiles {
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("read file [%s]: %w", filePath, err)
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
			return fmt.Errorf("scan lines in file [%s]: %w", filePath, err)
		}

		if len(lines) == 0 {
			continue
		}

		relPath, err := filepath.Rel(outputAbsPath, filePath)
		if err != nil {
			return fmt.Errorf("calculate relative path: %w", err)
		}

		// TODO: get line number of the comment
		link := fmt.Sprintf("source: %s:%d\n", relPath, 0)

		if _, err := outputFile.WriteString(link); err != nil {
			return fmt.Errorf("write line to output file: %w", err)
		}

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
