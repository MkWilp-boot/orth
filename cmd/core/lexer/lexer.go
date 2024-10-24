package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"orth/cmd/core/orth_debug"
	orth_types "orth/cmd/pkg/types"
	"os"
	"path"
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

func preProccessFile(includeFile string, parsedFiles chan orth_types.File[string]) {
	file, _ := os.Open(includeFile)
	// looks up for a file with same name on the include paths provided by the programmer
	if file == nil && *orth_debug.I != "" {
		paths := strings.Split(*orth_debug.I, ",")
		for i := 0; i < len(paths); i++ {
			paths[i] = strings.Trim(paths[i], " ")
			paths[i] = path.Join(paths[i], includeFile)

			if _, err := os.Stat(paths[i]); errors.Is(err, os.ErrNotExist) {
				continue
			}
			file, _ = os.Open(paths[i])
		}
	}
	if file == nil {
		fmt.Fprint(os.Stderr, orth_debug.BuildErrorMessage(orth_debug.ORTH_ERR_15, includeFile))
		os.Exit(1)
	}

	defer file.Close()

	var rawFile string
	source := orth_types.File[string]{
		Name:      includeFile,
		CodeBlock: rawFile,
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		rawFile = fmt.Sprintf("%s %s", rawFile, line)
		source.UpdateCodeReference(rawFile)

		if len(line) <= 0 || !strings.HasPrefix(line, "@") {
			continue
		}
		var directive string

		for i := 1; i < len(line) && line[i] != ' '; i++ {
			directive += string(line[i])
		}

		switch directive {
		case "define":
			name, value := ppDefineDirective(line)
			rawFile = strings.Replace(rawFile, fmt.Sprintf("@define %s %s", name, value), "", -1)
			rawFile = strings.ReplaceAll(rawFile, name, value)

			source.UpdateCodeReference(rawFile)
		case "include":
			includeFile := ""
			for i := len(directive) + 2; i < len(line) && line[i] != ' '; i++ {
				includeFile += string(line[i])
			}
			includeFile = strings.ReplaceAll(includeFile, `"`, "")
			includeFile = strings.TrimSpace(includeFile)

			rawFile = strings.Replace(rawFile, fmt.Sprintf(`@include "%s"`, includeFile), "", -1)
			source.UpdateCodeReference(rawFile)

			filesToParse := make(chan orth_types.File[string])

			go preProccessFile(includeFile, filesToParse)

			for file := range filesToParse {
				parsedFiles <- file
			}
		default:
			fmt.Printf("unknow directive found: %q is not recognized as an internal or external directive\n", directive)
			os.Exit(2)
		}
	}

	parsedFiles <- source

	close(parsedFiles)
}

func LoadProgramFromFile(path string) []orth_types.File[string] {
	// removePathToFile := regexp.MustCompile(`((\.\.\/|\.\/)+|("))`)
	// path = removePathToFile.ReplaceAllString(path, "")
	filesParsed := make(chan orth_types.File[string])
	go preProccessFile(path, filesParsed)

	files := make([]orth_types.File[string], 0)

	for file := range filesParsed {
		files = append(files, file)
	}

	return files
}

// LexFile receives a pure text program then
// separate and enumerates all tokens present within the provided program
func LexFile(programFiles []orth_types.File[string]) []orth_types.File[orth_types.SliceOf[orth_types.StringEnum]] {
	lexedFiles := make([]orth_types.File[orth_types.SliceOf[orth_types.StringEnum]], 0)

	for _, file := range programFiles {
		pLines := strings.Split(file.CodeBlock, "\r\n")
		lines := make([]orth_types.StringEnum, 0)

		for lineNumber, line := range pLines {
			if len(line) == 0 {
				continue
			}
			enumeration := make(chan orth_types.Vec2DString)

			go EnumerateLine(line, enumeration)

			for enumeratedLine := range enumeration {
				vec2d := orth_types.StringEnum{
					Index:   lineNumber + 1,
					Content: enumeratedLine,
				}
				lines = append(lines, vec2d)
			}
		}

		lexedFiles = append(lexedFiles, orth_types.File[orth_types.SliceOf[orth_types.StringEnum]]{
			Name: file.Name,
			CodeBlock: orth_types.SliceOf[orth_types.StringEnum]{
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
func EnumerateLine(line string, enumeration chan<- orth_types.Vec2DString) {
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

		enumeration <- orth_types.Vec2DString{
			Index: col,
			Token: line[col:colEnd],
		}

		col = findCol(line, colEnd, func(s string) bool {
			return s != " "
		})
	}
	close(enumeration)
}
