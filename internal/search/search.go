package search

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Options struct {
	IgnoreCase            bool
	ShowLineNumbers       bool
	PrintCountPerFile     bool
	PrintFilesWithMatches bool
}

type Match struct {
	LineNumber int
	Line       string
}

func (m Match) Format(showLineNumbers bool) string {
	if showLineNumbers {
		return fmt.Sprintf("%d: %s", m.LineNumber, m.Line)
	}
	return m.Line
}

func SearchFile(path, pattern string, opts Options) ([]Match, error) {
	// open file from file path and handle error
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	// close file after function runs
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lineNumber := 0
	var matches []Match

	// scan the file, and return the lines that match the pattern
	for scanner.Scan() {
		lineNumber++

		var patternMatches bool
		// lower cases if opts.IgnoreCase is true
		if opts.IgnoreCase {
			patternMatches = strings.Contains(strings.ToLower(scanner.Text()), strings.ToLower(pattern))
		} else {
			patternMatches = strings.Contains(scanner.Text(), pattern)
		}

		if patternMatches {
			matches = append(matches, Match{
				LineNumber: lineNumber,
				Line:       scanner.Text(),
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

func SearchFileForMatchCount(path, pattern string, opts Options) (string, error) {
	matches, err := SearchFile(path, pattern, opts)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", len(matches)), nil
}

func SearchFileForMatches(path, pattern string, opts Options) (bool, error) {
	matches, err := SearchFile(path, pattern, opts)
	if err != nil {
		return false, err
	}

	return len(matches) > 0, nil
}
