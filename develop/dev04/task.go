package main

import (
	"fmt"
	"sort"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак''пятка' и 'тяпка' - принадлежат одному множеству,
'листок''слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массивкаждый элемент которогослово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func hash(s string) string {
	h := make([]int, 33)
	for _, r := range s {
		alphPosition := r - 'а'
		if r == 'ё' {
			alphPosition = 32
		}
		h[alphPosition]++
	}
	return fmt.Sprint(h)
}

// AnagramSet .
func AnagramSet(arr []string) *map[string][]string {
	mapSet := make(map[string][]string)

	for _, word := range arr {
		word = strings.ToLower(word)
		h := hash(word)
		mapSet[h] = append(mapSet[h], word)
	}

	result := make(map[string][]string, len(mapSet))
	for _, set := range mapSet {
		if len(set) > 1 {
			key := set[0]
			sort.Strings(set)
			result[key] = set
		}
	}

	return &result
}

func main() {
	dict := []string{"тяпка", "пятка"}
	fmt.Println(AnagramSet(dict))
}
