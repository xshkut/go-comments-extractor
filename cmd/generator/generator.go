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

	fileStream := NewStream[string](64)
	lineStream := NewStream[[]string](64)

	go g.getGoFiles(inputAbsPath, fileStream)

	outputFile, err := os.Create(g.outputPath)
	if err != nil {
		return fmt.Errorf("create output file [%s]: %w", g.outputPath, err)
	}
	defer outputFile.Close()

	go g.extractCommentsFromFilesAndSave(fileStream, lineStream, outputAbsPath, outputFile)

	err = g.saveLines(lineStream, outputAbsPath, outputFile)
	if err != nil {
		return fmt.Errorf("extract comments: %w", err)
	}

	return nil
}

func (g *commentsExtractor) extractCommentsFromFilesAndSave(fileStream *streamChain[string], lineStream *streamChain[[]string], outputAbsPath string, outputFile *os.File) {
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

func (g *commentsExtractor) saveLines(outpLines *streamChain[[]string], outputAbsPath string, outputFile *os.File) error {
	var err error

	if g.header != "" {
		headerContent := []string{g.header, "\n"}

		err = appendContent(outputFile, headerContent)
		if err != nil {
			return fmt.Errorf("append header: %w", err)
		}
	}

	for comments := range outpLines.Chan() {
		if err := outpLines.Error(); err != nil {
			return err
		}

		err = appendContent(outputFile, comments)
		if err != nil {
			return fmt.Errorf("append content: %w", err)
		}
	}

	return nil
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

func (g *commentsExtractor) getGoFiles(dir string, st *streamChain[string]) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walking to %s: %w", path, err)
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			st.Chan() <- path
		}

		return nil
	})

	if err != nil {
		st.Destroy(fmt.Errorf("walk over files:%w", err))

		return
	}

	st.End()
}
