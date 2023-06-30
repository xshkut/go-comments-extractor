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

	filePaths := newStream[string](64)
	content := newStream[[]string](64)
	result := newStream[struct{}](64)

	go g.getGoFiles(inputAbsPath, filePaths)
	go g.extractComments(filePaths, content, outputAbsPath)
	go g.saveContent(content, result, outputAbsPath)

	<-result.Chan()

	if err := result.Error(); err != nil {
		return fmt.Errorf("chain error: %w", err)
	}

	return nil
}

func (g *commentsExtractor) getGoFiles(dir string, files output[string]) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walking to %s: %w", path, err)
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			files.Write(path)
		}

		return nil
	})

	if err != nil {
		files.Destroy(fmt.Errorf("walk over files:%w", err))

		return
	}

	files.End()
}

func (g *commentsExtractor) extractComments(fileStream input[string], lineStream output[[]string], outputAbsPath string) {
	i := 0

	for filePath := range fileStream.Chan() {
		if err := fileStream.Error(); err != nil {
			lineStream.Destroy(err)
		}

		i++

		relPathLink, err := getRelPathLink(outputAbsPath, filePath)
		if err != nil {
			lineStream.Destroy(fmt.Errorf("calculate relative path: %w", err))
			return
		}

		link := fmt.Sprintf("%s source: %s", g.outputCommentPrefix, relPathLink)

		content, err := os.ReadFile(filePath)
		if err != nil {
			lineStream.Destroy(fmt.Errorf("read file [%s]: %w", filePath, err))
			return
		}

		lines, err := scanLines(content, g.inputPattern)
		if err != nil {
			lineStream.Destroy(fmt.Errorf("scan lines in file [%s]: %w", filePath, err))
			return
		}

		if len(lines) == 0 {
			continue
		}

		lines = append(lines, "")
		lines = append([]string{link, ""}, lines...)

		lineStream.Write(lines)
	}

	lineStream.End()
}

func (g *commentsExtractor) saveContent(outpLines input[[]string], result output[struct{}], outputAbsPath string) {
	outputFile, err := os.Create(g.outputPath)
	if err != nil {
		result.Destroy(fmt.Errorf("create output file [%s]: %w", g.outputPath, err))
		return
	}
	defer outputFile.Close()

	if g.header != "" {
		headerContent := []string{g.header, "\n"}

		err = appendContent(outputFile, headerContent)
		if err != nil {
			result.Destroy(fmt.Errorf("append header: %w", err))
			return
		}
	}

	for comments := range outpLines.Chan() {
		if err := outpLines.Error(); err != nil {
			result.Destroy(err)
			return
		}

		err = appendContent(outputFile, comments)
		if err != nil {
			result.Destroy(fmt.Errorf("append content: %w", err))
			return
		}
	}

	result.End()
}

func getRelPathLink(outputAbsPath string, filePath string) (string, error) {
	outputFolder := filepath.Dir(outputAbsPath)
	relPath, err := filepath.Rel(outputFolder, filePath)

	relPath = "file://./" + relPath

	return relPath, err
}

func appendContent(outputFile *os.File, lines []string) error {
	schemaContent := strings.Join(lines, "\n")

	if _, err := outputFile.WriteString(schemaContent); err != nil {
		return fmt.Errorf("write line to output file: %w", err)
	}

	return nil
}

func scanLines(content []byte, inputPattern string) ([]string, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	scanner.Split(bufio.ScanLines)

	lines := make([]string, 0)

	var extracting bool

	multilinePrefix := fmt.Sprintf("/* %s:", inputPattern)
	singlelinePrefix := fmt.Sprintf("// %s:", inputPattern)

	for scanner.Scan() {
		line := scanner.Text()

		if extracting {
			if strings.HasPrefix(line, "*/") {
				extracting = false
				lines = append(lines, "")
				continue
			}

			lines = append(lines, line)
		} else if strings.HasPrefix(line, multilinePrefix) {
			extracting = true
		} else if strings.HasPrefix(line, singlelinePrefix) {
			line = strings.TrimPrefix(line, singlelinePrefix)
			line = strings.TrimSpace(line)

			lines = append(lines, line)
			lines = append(lines, "")
		} else {
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return lines, nil
}
