package lexer

import (
	"io/ioutil"
	"strings"
	orthtypes "t/cmd/pkg/types"
)

func LoadProgramFromFile(path string) []orthtypes.StringEnum {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return LexFile(string(fileBytes))
}

func LexFile(strProgram string) []orthtypes.StringEnum {
	lines := make([]orthtypes.StringEnum, 0)
	for lineNumber, line := range strings.Split(strProgram, "\r\n") {
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
	return lines
}

func findCol(line string, start int, predicate func(string) bool) int {
	for start < len(line) && !predicate(string(line[start])) {
		start++
	}
	return start
}

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
