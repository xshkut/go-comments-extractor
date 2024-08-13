package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type extractor struct {
	inputPath           string
	outputPath          string
	inputPattern        string
	outputCommentPrefix string
	header              string
}

type CommentsExtractorConfig struct {
	InputPath           string
	OutputPath          string
	InputCommentPattern string
	OutputCommentPrefix string
	OutputFileHeader    string
}

func NewCommentsExtractor(cfg CommentsExtractorConfig) *extractor {
	return &extractor{
		inputPath:           cfg.InputPath,
		outputPath:          cfg.OutputPath,
		inputPattern:        cfg.InputCommentPattern,
		outputCommentPrefix: cfg.OutputCommentPrefix,
		header:              cfg.OutputFileHeader,
	}
}

func (g *extractor) ExtractComments() error {
	inputAbsPath, err := filepath.Abs(g.inputPath)
	if err != nil {
		return fmt.Errorf("get abs path [%s]: %w", g.inputPath, err)
	}

	outputAbsPath, err := filepath.Abs(g.outputPath)
	if err != nil {
		return fmt.Errorf("get abs path [%s]: %w", g.outputPath, err)
	}

	ctx := context.Context(context.Background())

	filePaths := newStream[string](64)
	content := newStream[[]string](64)
	result := newStream[struct{}](0)

	go g.getGoFiles(ctx, inputAbsPath, filePaths)
	go g.extractComments(filePaths, content, outputAbsPath)
	go g.saveContent(content, result, outputAbsPath)

	<-result.Chan()

	ctx.Done()

	return wrapIfError(err, "chain error")
}

func (g *extractor) getGoFiles(ctx context.Context, dir string, files output[string]) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("context error: %w", err)
		}

		if err != nil {
			return fmt.Errorf("walking to %s: %w", path, err)
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			files.Write(path)
		}

		return nil
	})

	files.End(wrapIfError(err, "walk over files"))
}

func (g *extractor) extractComments(fileStream input[string], lineStream output[[]string], outputAbsPath string) {
	i := 0

	for filePath := range fileStream.Chan() {
		i++

		relPathLink, err := getRelPathLink(outputAbsPath, filePath)
		if err != nil {
			lineStream.End(fmt.Errorf("calculate relative path: %w", err))
			return
		}

		link := fmt.Sprintf("%s source: %s", g.outputCommentPrefix, relPathLink)

		content, err := os.ReadFile(filePath)
		if err != nil {
			lineStream.End(fmt.Errorf("read file [%s]: %w", filePath, err))
			return
		}

		lines, err := scanLines(content, g.inputPattern)
		if err != nil {
			lineStream.End(fmt.Errorf("scan lines in file [%s]: %w", filePath, err))
			return
		}

		if len(lines) == 0 {
			continue
		}

		lines = append(lines, "")
		lines = append([]string{link, ""}, lines...)

		lineStream.Write(lines)
	}

	lineStream.End(fileStream.Error())
}

func (g *extractor) saveContent(outpLines input[[]string], result output[struct{}], outputAbsPath string) {
	outputFile, err := os.Create(outputAbsPath)
	if err != nil {
		result.End(fmt.Errorf("create output file [%s]: %w", g.outputPath, err))
		return
	}
	defer outputFile.Close()

	if g.header != "" {
		headerContent := []string{g.header, "\n"}

		err = appendContent(outputFile, headerContent)
		if err != nil {
			result.End(fmt.Errorf("append header: %w", err))
			return
		}
	}

	for comments := range outpLines.Chan() {
		err = appendContent(outputFile, comments)
		if err != nil {
			result.End(fmt.Errorf("append content: %w", err))
			return
		}
	}

	result.End(outpLines.Error())
}

func getRelPathLink(outputAbsPath string, filePath string) (string, error) {
	outputFolder := filepath.Dir(outputAbsPath)
	relPath, err := filepath.Rel(outputFolder, filePath)

	relPath = "file://./" + relPath

	return relPath, err
}

func appendContent(outputFile *os.File, lines []string) error {
	schemaContent := strings.Join(lines, "\n")

	_, err := outputFile.WriteString(schemaContent)

	return wrapIfError(err, "write line to output file")
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

		if !extracting {
			line = strings.TrimSpace(line)
		}

		switch {
		case extracting && strings.HasPrefix(line, "*/"):
			extracting = false

			lines = append(lines, "")
		case extracting:
			lines = append(lines, line)
		case strings.HasPrefix(line, multilinePrefix):
			extracting = true
		case strings.HasPrefix(line, singlelinePrefix):
			line = strings.TrimPrefix(line, singlelinePrefix)
			line = strings.TrimSpace(line)

			lines = append(lines, line, "")
		}
	}

	return lines, wrapIfError(scanner.Err(), "scanner error")
}
