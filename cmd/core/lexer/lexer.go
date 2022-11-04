package lexer

import (
	"fmt"
	orthtypes "orth/cmd/pkg/types"
	"os"
	"regexp"
	"strings"
)

func getParams(regEx, line string) (paramsMap []string) {
	var compRegEx = regexp.MustCompile(regEx)
	match := compRegEx.FindAllStringSubmatch(line, -1)

	paramsMap = make([]string, 0)
	for i := range match {
		paramsMap = append(paramsMap, match[i][1])
	}
	return paramsMap
}

func getParamsMap(regEx, line string) (paramsMap []map[string]string) {
	var compRegEx = regexp.MustCompile(regEx)
	matchs := compRegEx.FindAllStringSubmatch(line, -1)
	paramsMap = make([]map[string]string, 0)

	for _, match := range matchs {
		matchMap := make(map[string]string)
		matchMap[match[1]] = match[2]
		paramsMap = append(paramsMap, matchMap)
	}
	return paramsMap
}

func LoadProgramFromFile(path string) []orthtypes.File[string] {
	fileBytes, err := os.ReadFile(path)
	removePathToFile := regexp.MustCompile(`((\.\.\/|\.\/)+|("))`)
	path = removePathToFile.ReplaceAllString(path, "")

	if err != nil {
		panic(err)
	}

	strProgram := string(fileBytes)

	defineDirective := getParamsMap(`@define\s(?P<DName>\w+)\s(?P<Replacement>[\w #$%*\(\)\\\/"]*)`, strProgram)

	for _, s := range defineDirective {
		for k, v := range s {
			str := fmt.Sprintf("@define %s %s", k, v)
			strProgram = strings.ReplaceAll(strProgram, str, "")
			strProgram = strings.ReplaceAll(strProgram, k, v)
		}
	}

	files := make([]orthtypes.File[string], 1)
	files[0] = orthtypes.File[string]{
		Name:      path,
		CodeBlock: strProgram,
	}

	includeFiles := getParams(`(?i)@include\s"(?P<File>\w+\.orth)"`, strProgram)

	for _, v := range includeFiles {
		rmInclude := regexp.MustCompile(`(?i)@include\s"` + v + `"\r?\n?`)
		strProgram = rmInclude.ReplaceAllString(strProgram, "")
		files[0].CodeBlock = strProgram

		includedProgram := LoadProgramFromFile(v)
		files = append(files, includedProgram...)
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
	line = strings.Split(line, "//")[0]
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
			Index:   col,
			Content: line[col:colEnd],
		}

		col = findCol(line, colEnd, func(s string) bool {
			return s != " "
		})
	}
	close(enumeration)
}
