package main

import (
	"reflect"
	"testing"
)

func TestAnagramSet(t *testing.T) {
	testTable := []struct {
		input  []string
		result *map[string][]string
	}{
		{
			input: []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик"},
			result: &map[string][]string{
				"пятак": []string{
					"пятак",
					"пятка",
					"тяпка",
				},
				"листок": []string{
					"листок",
					"слиток",
					"столик",
				},
			},
		},
		{
			input: []string{"тяпка", "листок", "слиток", "столик"},
			result: &map[string][]string{
				"листок": []string{
					"листок",
					"слиток",
					"столик",
				},
			},
		},
		{
			input: []string{"тяпка", "слиток", "листок", "столик"},
			result: &map[string][]string{
				"слиток": []string{
					"листок",
					"слиток",
					"столик",
				},
			},
		},
		{
			input:  []string{"тяпка", "столик"},
			result: &map[string][]string{},
		},
		{
			input:  []string{},
			result: &map[string][]string{},
		},
		{
			input:  []string{"б"},
			result: &map[string][]string{},
		},
	}

	for _, testCase := range testTable {
		result := AnagramSet(testCase.input)

		t.Logf("Calling anagramSet(%v), result %v", testCase.input, result)

		if !reflect.DeepEqual(*result, *testCase.result) {
			t.Errorf("Incorrect result: expect %v, got %v",
				testCase.result, result)
		}
	}
}

func TestHash(t *testing.T) {
	testTable := []struct {
		input  string
		result string
	}{
		{
			input:  "абвгд",
			result: "[1 1 1 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]",
		},
		{
			input:  "е",
			result: "[0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]",
		},
		{
			input:  "ё",
			result: "[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1]",
		},
		{
			input:  "",
			result: "[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]",
		},
	}

	for _, testCase := range testTable {
		result := hash(testCase.input)

		t.Logf("Calling hash(%s), result %s", testCase.input, result)

		if result != testCase.result {
			t.Errorf("Incorrect result: expect %s, got %s", testCase.result, result)
		}
	}
}
