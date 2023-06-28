# Go Comments Extractor

Extracts comments marked with a defined prefix into a single file.

## Installation

```bash
go get github.com/xshkut/go-comments-extractor
```

## Usage

```bash
go run github.com/xshkut/go-comments-extractor/cmd/generator -i ./playground/go -o ./playground/schema.sql -p SQL -c "--" -h "Generated SQL Schema"
```

```
-c string
      Prefix for comments in the output file. Optional
-h string
      Header of the output file. Optional
-i string
      Root input path (file or folder) (default "./")
-o string
      Output file. Required
-p string
      Prefix going right after the beginning of Go multiline comment plus whitespace. Required
```

## Motivation:

1. AI code generation: In most cases, AI code generation and autocompletion works better when the referenced source code is placed in the same file where we want to autocomplete. This tools addresses this case.
2. In-place documentation: put docs near related entities and generate aggregated docs in a single file.

## Example:

In [example/go](./example/go) we have sql schema distributed across Go source code for coupled referrencing and AI autocompletion.
After running the tool (command below) we will have a single file [example/schema.sql](./example/schema.sql) with all SQL schema components. Later that may be used to see the schema diffs in a single place.

```bash
go run github.com/xshkut/go-comments-extractor/cmd/generator -i ./playground/go -o ./playground/schema.sql -p SQL -c "--" -h "Generated SQL Schema"
```