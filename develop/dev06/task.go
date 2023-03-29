package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var (
	fields    string
	delimiter string
	separated bool

	fileNames []string
)

var (
	errorFieldsFlag    = errors.New("The -f flag is required")
	errorFileNotFound  = errors.New("No such file or directory")
	errorInvalidRegexp = errors.New("Invalid regular expression")
	errorNoFiles       = errors.New("File to cut was not specified")
)

func init() {
	flag.StringVar(&fields, "f", "", "select fields")
	flag.StringVar(&delimiter, "d", "\t", "set delimiter")
	flag.BoolVar(&separated, "s", false, "delimited lines only")

	flag.Parse()
	fileNames = flag.Args()
}

func main() {
	if fields == "" {
		fmt.Println(errorFieldsFlag)
		return
	}

	if len(fileNames) == 0 {
		fmt.Println(errorNoFiles)
		return
	}

	flds, err := validateFields(fields)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, fileName := range fileNames {
		input, err := os.Open(fileName)
		if err != nil {
			fmt.Printf("%s: %s\n", errorFileNotFound, fileName)
		}

		err = cut(input, flds)
		if err != nil {
			fmt.Println(err)
		}
		input.Close()
	}
}

func validateFields(fields string) ([]int, error) {
	var flds []int

	for _, fld := range strings.Split(fields, ",") {
		fld = strings.TrimSpace(fld)

		v, err := strconv.Atoi(fld)
		if err != nil {
			return nil, fmt.Errorf("Error in validateFields - strconv.Atoi with -f %s: %s", fields, err.Error())
		}
		flds = append(flds, v-1)
	}

	return flds, nil
}

func cut(file *os.File, flds []int) error {
	buffer := bytes.Buffer{}
	byteDel := []byte(delimiter)

	scanner := bufio.NewScanner(file)

	firstLine := true
	for scanner.Scan() {
		line := scanner.Bytes()

		isContain := bytes.Contains(line, byteDel)
		if separated && !isContain {
			continue
		}

		if !firstLine {
			buffer.WriteRune('\n')
		}
		firstLine = false

		if !isContain {
			buffer.Write(line)
			continue
		}

		words := bytes.Split(line, byteDel)

		first := true
		for _, fld := range flds {
			if fld < len(words) {
				if !first {
					buffer.Write(byteDel)
				}
				buffer.Write(words[fld])
				first = false
			}
		}
	}

	fmt.Println(buffer.String())

	return nil
}
