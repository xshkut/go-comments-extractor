package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type commentsExtractor struct {
	inputPath           string
	outputPath          string
	inputPattern        string
	outputCommentPrefix string
	header              string
}

func (g *commentsExtractor) ExtractComments() error {
	inputAbsPath, err := filepath.Abs(g.inputPath)
	if err != nil {
		return fmt.Errorf("get abs path [%s]: %w", g.inputPath, err)
	}

	outputAbsPath, err := filepath.Abs(g.outputPath)
	if err != nil {
		return fmt.Errorf("get abs path [%s]: %w", g.outputPath, err)
	}

	goFiles, err := g.getGoFiles(inputAbsPath)
	if err != nil {
		return fmt.Errorf("list files [%s]: %w", inputAbsPath, err)
	}

	outputFile, err := os.Create(g.outputPath)
	if err != nil {
		return fmt.Errorf("create output file [%s]: %w", g.outputPath, err)
	}
	defer outputFile.Close()

	err = g.extractCommentsFromFiles(goFiles, outputAbsPath, outputFile)
	if err != nil {
		return fmt.Errorf("extract comments: %w", err)
	}

	return nil
}

func (g *commentsExtractor) extractCommentsFromFiles(goFiles []string, outputAbsPath string, outputFile *os.File) error {
	prependWithNewLine := false

	for _, filePath := range goFiles {
		relPath, err := filepath.Rel(outputAbsPath, filePath)
		if err != nil {
			return fmt.Errorf("calculate relative path: %w", err)
		}

		link := fmt.Sprintf("%s source: %s\n", g.outputCommentPrefix, relPath)

		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("read file [%s]: %w", filePath, err)
		}

		lines, err := scanLines(content, g)
		if err != nil {
			return fmt.Errorf("scan lines in file [%s]: %w", filePath, err)
		}

		if len(lines) == 0 {
			continue
		}

		err = appendContent(prependWithNewLine, outputFile, link, lines)
		if err != nil {
			return fmt.Errorf("append content: %w", err)
		}

		prependWithNewLine = true
	}

	return nil
}

func appendContent(prependWithNewLine bool, outputFile *os.File, link string, lines []string) error {
	if prependWithNewLine {
		if _, err := outputFile.WriteString("\n"); err != nil {
			return fmt.Errorf("write line to output file: %w", err)
		}
	}

	if _, err := outputFile.WriteString(link); err != nil {
		return fmt.Errorf("write line to output file: %w", err)
	}

	schemaContent := strings.Join(lines, "\n")

	if _, err := outputFile.WriteString(schemaContent); err != nil {
		return fmt.Errorf("write line to output file: %w", err)
	}

	if _, err := outputFile.WriteString("\n"); err != nil {
		return fmt.Errorf("write line to output file: %w", err)
	}

	return nil
}

func scanLines(content []byte, g *commentsExtractor) ([]string, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	scanner.Split(bufio.ScanLines)

	lines := make([]string, 0)

	var extracting bool

	for scanner.Scan() {
		line := scanner.Text()

		if extracting {
			if strings.HasPrefix(line, "*/") {
				extracting = false
				break
			}

			lines = append(lines, line)
		} else if strings.HasPrefix(line, fmt.Sprintf("/* %s", g.inputPattern)) {
			extracting = true
			lines = append(lines, strings.TrimPrefix(line, fmt.Sprintf("/* %s", g.inputPattern)))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return lines, nil
}

func (g *commentsExtractor) getGoFiles(dir string) ([]string, error) {
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
