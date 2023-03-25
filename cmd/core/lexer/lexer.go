package lexer

import (
	"fmt"
	orthtypes "orth/cmd/pkg/types"
	"os"
	"regexp"
	"strings"
)

func ppDefineDirective(line string) (string, string) {
	const directive = "define"
	name := ""
	for i := len(directive) + 2; i < len(line) && line[i] != ' '; i++ {
		name += string(line[i])
	}
	value := ""
	for i := len(directive) + len(name) + 3; i < len(line); i++ {
		value += string(line[i])
	}
	return strings.TrimSpace(name), strings.TrimSpace(value)
}

func proProccessFile(rawFile, path string, parsedFiles chan orthtypes.File[string]) {
	lines := strings.Split(rawFile, "\n")

	oFile := orthtypes.File[string]{
		Name:      path,
		CodeBlock: rawFile,
	}

	for _, line := range lines {
		if len(line) <= 0 || !strings.HasPrefix(line, "@") {
			continue
		}

		directive := ""

		for i := 1; i < len(line) && line[i] != ' '; i++ {
			directive += string(line[i])
		}

		switch directive {
		case "define":
			name, value := ppDefineDirective(line)
			rawFile = strings.Replace(rawFile, fmt.Sprintf("@define %s %s", name, value), "", -1)
			rawFile = strings.ReplaceAll(rawFile, name, value)

			oFile.UpdateCodeReference(rawFile)
		case "include":
			includeFile := ""
			for i := len(directive) + 2; i < len(line) && line[i] != ' '; i++ {
				includeFile += string(line[i])
			}
			includeFile = strings.ReplaceAll(includeFile, `"`, "")
			includeFile = strings.TrimSpace(includeFile)

			includeFileContent, err := os.ReadFile(includeFile)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			rawFile = strings.Replace(rawFile, fmt.Sprintf(`@include "%s"`, includeFile), "", -1)
			oFile.UpdateCodeReference(rawFile)

			filesToParse := make(chan orthtypes.File[string])

			go proProccessFile(string(includeFileContent), includeFile, filesToParse)

			for file := range filesToParse {
				parsedFiles <- file
			}
		default:
			fmt.Printf("unknow directive found: %q is not recognized as an internal or external directive\n", directive)
			os.Exit(2)
		}
	}

	parsedFiles <- oFile

	close(parsedFiles)
}

func LoadProgramFromFile(path string) []orthtypes.File[string] {
	fileBytes, err := os.ReadFile(path)
	removePathToFile := regexp.MustCompile(`((\.\.\/|\.\/)+|("))`)
	path = removePathToFile.ReplaceAllString(path, "")

	if err != nil {
		panic(err)
	}

	filesParsed := make(chan orthtypes.File[string])
	go proProccessFile(string(fileBytes), path, filesParsed)

	files := make([]orthtypes.File[string], 0)

	for file := range filesParsed {
		files = append(files, file)
	}

	return files
}

// LexFile receives a pure text program then
// separate and enumerates all tokens present within the provided program
func LexFile(programFiles []orthtypes.File[string]) []orthtypes.File[orthtypes.SliceOf[orthtypes.StringEnum]] {
	lexedFiles := make([]orthtypes.File[orthtypes.SliceOf[orthtypes.StringEnum]], 0)

	for _, file := range programFiles {
		pLines := strings.Split(file.CodeBlock, "\r\n")
		lines := make([]orthtypes.StringEnum, 0)

		for lineNumber, line := range pLines {
			if len(line) == 0 {
				continue
			}
			enumeration := make(chan orthtypes.Vec2DString)

			go EnumerateLine(line, enumeration)

			for enumeratedLine := range enumeration {
				vec2d := orthtypes.StringEnum{
					Index:   lineNumber + 1,
					Content: enumeratedLine,
				}
				lines = append(lines, vec2d)
			}
		}

		lexedFiles = append(lexedFiles, orthtypes.File[orthtypes.SliceOf[orthtypes.StringEnum]]{
			Name: file.Name,
			CodeBlock: orthtypes.SliceOf[orthtypes.StringEnum]{
				Slice: &lines,
			},
		})
	}
	return lexedFiles
}

// findCol separates the tokens in a `line` starting at `start` by executing a predicate
func findCol(line string, start int, predicate func(string) bool) int {
	for start < len(line) && !predicate(string(line[start])) {
		start++
	}
	return start
}

// EnumerateLine receives a single line and parses and enumerates
// all tokens in that line feeding the `enumeration` chan
func EnumerateLine(line string, enumeration chan<- orthtypes.Vec2DString) {
	line = strings.Split(line, "#")[0]
	col := findCol(line, 0, func(s string) bool {
		return s != " "
	})

	continueSearch := true

	for col < len(line) {
		colEnd := findCol(line, col, func(s string) bool {
			if s == "\"" {
				continueSearch = !continueSearch
			}
			if !continueSearch {
				return false
			}
			return s == " "
		})

		enumeration <- orthtypes.Vec2DString{
			Index: col,
			Token: line[col:colEnd],
		}

		col = findCol(line, colEnd, func(s string) bool {
			return s != " "
		})
	}
	close(enumeration)
}
