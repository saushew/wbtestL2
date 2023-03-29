package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var (
	after      int
	before     int
	context    int
	count      bool
	ignoreCase bool
	invert     bool
	fixed      bool
	lineNum    bool

	searchWord string
	fileNames  []string
)

var (
	errorNotEnoughArgs = errors.New("Not enough args")
	errorFileNotFound  = errors.New("No such file or directory")
	errorInvalidRegexp = errors.New("Invalid regular expression")
)

func init() {
	flag.IntVar(&after, "A", 0, "print +N lines after a match")
	flag.IntVar(&before, "B", 0, "print +N lines before a match")
	flag.IntVar(&context, "C", 0, "print ±N lines around the match")
	flag.BoolVar(&count, "c", false, "print only a count of selected lines")
	flag.BoolVar(&ignoreCase, "i", false, "ignore case differences")
	flag.BoolVar(&invert, "v", false, "select non-matching lines")
	flag.BoolVar(&fixed, "F", false, "exact match with a string, not a pattern")
	flag.BoolVar(&lineNum, "n", false, "sort by numeric value, taking into account suffixes")

	flag.Parse()
	searchWord = flag.Arg(0)
	fileNames = flag.Args()
}

func main() {
	if len(fileNames) < 2 {
		fmt.Println(errorNotEnoughArgs)
		return
	}
	fileNames = fileNames[1:]

	if context > 0 {
		after = context
		before = context
	}

	if ignoreCase {
		searchWord = strings.ToLower(searchWord)
	}

	var check func(line string) bool
	if fixed {
		check = func(line string) bool {
			return strings.Contains(line, searchWord) != invert
		}
	} else {
		re, err := regexp.Compile(searchWord)
		if err != nil {
			fmt.Printf("%s: %s: %s\n", searchWord, errorInvalidRegexp.Error(), err.Error())
			return
		}

		check = func(line string) bool {
			return re.MatchString(line) != invert
		}
	}

	for _, fileName := range fileNames {
		err := search(fileName, check)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func search(fileName string, check func(string) bool) error {
	lines, err := readFile(fileName)
	if err != nil {
		return err
	}

	builder := strings.Builder{}
	var counter int
	for i, line := range lines {
		if ignoreCase {
			line = strings.ToLower(line)
		}

		if !check(line) {
			continue
		}

		if count {
			counter++
			continue
		}

		from, to := i-before, i+after+1
		if from < 0 {
			from = 0
		}
		if to > len(lines) {
			to = len(lines)
		}

		for k := from; k < to; k++ {
			if builder.Len() > 0 {
				builder.WriteString("\n")
			}
			if lineNum {
				builder.WriteString(fmt.Sprintf("%d:", i))
			}
			builder.WriteString(lines[k])
		}
	}

	if count {
		fmt.Println(counter)
		return nil
	}

	if builder.Len() > 0 {
		fmt.Println(builder.String())
	}

	return nil
}

func readFile(fileName string) ([]string, error) {
	input, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errorFileNotFound, fileName)
	}
	b, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, fmt.Errorf("Error with file %s in search - readFile - ioutil.ReadAll(): %w", fileName, err)
	}
	return strings.Split(string(b), "\n"), nil
}
