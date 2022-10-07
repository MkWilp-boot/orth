package lexer

import (
	"io/ioutil"
	orthtypes "orth/cmd/pkg/types"
	"strings"
)

// LoadProgramFromFile receives a path for a program and returns LexFile(path)
func LoadProgramFromFile(path string) []orthtypes.StringEnum {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return LexFile(string(fileBytes))
}

// LexFile receives a pure text program then
// separate and enumerates all tokens present within the provided program
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
