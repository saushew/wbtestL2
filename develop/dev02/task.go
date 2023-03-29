package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var (
	errorUnpack = errors.New("invalid input")
)

func lastIndex(r []rune) int {
	for i := len(r) - 1; i >= 0; i-- {
		if !unicode.IsDigit(r[i]) {
			return i
		}
	}
	return -1
}

func reverse(s string) string {
	r := []rune(s)
	l := len(r)

	for i, j := 0, l-1; i < l/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}

	return string(r)
}

// UnpackString .
func UnpackString(input string) (string, error) {
	runeInput := []rune(input)

	var result strings.Builder
	for len(runeInput) > 0 {
		letterIndex := lastIndex(runeInput)
		if letterIndex == -1 {
			return "", errorUnpack
		}

		var slashFlag bool
		if runeInput[letterIndex] == '\\' {
			if letterIndex > 0 && runeInput[letterIndex-1] == '\\' {
				slashFlag = true
			} else if letterIndex < len(runeInput)-1 {
				slashFlag = true
				letterIndex++
			} else {
				return "", errorUnpack
			}
		}

		multiply := 1
		if letterIndex < len(runeInput)-1 {
			multiply, _ = strconv.Atoi(string(runeInput[letterIndex+1:]))
		}

		result.WriteString(strings.Repeat(string(runeInput[letterIndex]), multiply))
		if slashFlag {
			runeInput = runeInput[:letterIndex-1]
		} else {
			runeInput = runeInput[:letterIndex]
		}
	}

	return reverse(result.String()), nil
}

func main() {
	s := string([]byte{'q', 'w', '\\'})
	fmt.Println(s)
	result, err := UnpackString(s)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
}
