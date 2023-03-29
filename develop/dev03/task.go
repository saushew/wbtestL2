package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var (
	column       int
	numericValue bool
	reverse      bool
	unique       bool
	monthName    bool
	tailSpaces   bool
	sorted       bool
	suffix       bool

	fileNames []string
)

var (
	errorNoFiles      = errors.New("No files to sort are provided in the command line args")
	errorFileNotFound = errors.New("No such file or directory")
)

// Sorter .
type Sorter struct {
	data []string
	sort func([]string) int
	ptr  int
}

func init() {
	flag.IntVar(&column, "k", 0, "num of column to sort")
	flag.BoolVar(&numericValue, "n", false, "sort by numeric value")
	flag.BoolVar(&reverse, "r", false, "sort in reverse order")
	flag.BoolVar(&unique, "u", false, "only unique strings")

	flag.BoolVar(&monthName, "M", false, "sort by month name")
	flag.BoolVar(&tailSpaces, "b", false, "ignore tail spaces")
	flag.BoolVar(&sorted, "c", false, "check if the data is sorted")
	flag.BoolVar(&suffix, "h", false, "sort by numeric value, taking into account suffixes")

	flag.Parse()
	fileNames = flag.Args()
}

func main() {
	if len(fileNames) == 0 {
		fmt.Println(errorNoFiles)
		return
	}

	srtr := NewSorter()

	for _, fileName := range fileNames {
		err := srtr.AddFile(fileName)
		if err != nil {
			fmt.Println(err)
			continue
		}

		srtr.Sort()

		fmt.Printf("%s:\n", fileName)
		for _, v := range srtr.data[:srtr.ptr] {
			fmt.Println(v)
		}
	}
}

// NewSorter .
func NewSorter() *Sorter {
	return &Sorter{
		sort: func(data []string) int {
			result := len(data)
			if (numericValue || monthName) && column == 0 {
				column = 1
			}

			if column <= 0 {
				sort.Slice(data, func(i, j int) bool {
					if tailSpaces {
						return strings.TrimSpace(data[i]) < strings.TrimSpace(data[j])
					}
					return data[i] < data[j]
				})
			} else {
				sort.Slice(data, func(i, j int) bool {
					datai := strings.Fields(data[i])
					dataj := strings.Fields(data[j])

					if len(datai) < column-1 && len(dataj) < column-1 {
						return data[i] < data[j]
					} else if len(datai) < column-1 {
						return true
					} else if len(dataj) < column-1 {
						return false
					}
					result, err := compare(datai[column-1], dataj[column-1])
					if err != nil {
						return data[i] < data[j]
					}
					return result
				})
			}

			if unique {
				ptr := 0
				for i := 1; i < len(data); i++ {
					if data[i] != data[ptr] {
						ptr++
						if ptr != i {
							data[ptr] = data[i]
						}
					}
				}
				result = ptr + 1
			}

			return result
		},
	}
}

func compare(a, b string) (bool, error) {
	switch {
	case numericValue:
		_a, errA := strconv.ParseFloat(a, 64)
		_b, errB := strconv.ParseFloat(b, 64)

		if errA != nil && errB != nil {
			return false, errA
		} else if errA != nil {
			return true, nil
		} else if errB != nil {
			return false, nil
		}
		return _a < _b, nil

	case monthName:
		_a, errA := time.Parse("Jan", a)
		_b, errB := time.Parse("Jan", b)

		if errA != nil && errB != nil {
			return false, errA
		} else if errA != nil {
			return true, nil
		} else if errB != nil {
			return false, nil
		}
		return _a.Month() < _b.Month(), nil
	default:
		return a < b, nil
	}
}

// Sort .
func (s *Sorter) Sort() {
	ptr := s.sort(s.data)
	s.ptr = ptr
}

// AddFile .
func (s *Sorter) AddFile(fileName string) error {
	s.data = nil

	data, err := Read(fileName)
	if err != nil {
		return err
	}

	s.data = data

	return nil
}

// Read .
func Read(fileName string) ([]string, error) {
	input, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errorFileNotFound, fileName)
	}
	b, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, fmt.Errorf("Error with file %s in Read - ioutil.ReadAll(): %w", fileName, err)
	}

	return strings.Split(string(b), "\n"), nil
}
