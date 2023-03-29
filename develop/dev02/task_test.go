package main

import "testing"

func TestUnpackString(t *testing.T) {
	testTable := []struct {
		input  string
		result string
		err    error
	}{
		{
			input:  "a4bc2d5e",
			result: "aaaabccddddde",
			err:    nil,
		},
		{
			input:  "abcd",
			result: "abcd",
			err:    nil,
		},
		{
			input:  "45",
			result: "",
			err:    errorUnpack,
		},
		{
			input:  "",
			result: "",
			err:    nil,
		},
		{
			input:  string([]byte{'q', 'w', 'e', '\\', '4', '\\', '5'}), // qwe\4\5
			result: "qwe45",
			err:    nil,
		},
		{
			input:  string([]byte{'q', 'w', 'e', '\\', '4', '5'}), // qwe\45
			result: "qwe44444",
			err:    nil,
		},
		{
			input:  string([]byte{'q', 'w', 'e', '\\', '\\', '5'}),              // qwe\\5
			result: string([]byte{'q', 'w', 'e', '\\', '\\', '\\', '\\', '\\'}), // qwe\\\\\
			err:    nil,
		},
		{
			input:  "4",
			result: "",
			err:    errorUnpack,
		},
		{
			input:  string([]byte{'\\', '5'}), // \5
			result: "5",
			err:    nil,
		},
		{
			input:  string([]byte{'\\', '\\'}), // \\
			result: string([]byte{'\\'}),       // \
			err:    nil,
		},
		{
			input:  string([]byte{'\\'}), // \
			result: "",       
			err:    errorUnpack,
		},
		{
			input:  string([]byte{'q', 'w', 'e', '\\'}), // qwe\
			result: "",       
			err:    errorUnpack,
		},
	}

	for _, testCase := range testTable {
		result, err := unpackString(testCase.input)

		t.Logf("Calling unpackString(%s), result %s, error %v", testCase.input, result, err)

		if result != testCase.result || err != testCase.err {
			t.Errorf("Incorrect result: expect (%s, %s), got (%s, %s)",
				testCase.result, testCase.err,
				result, err)
		}
	}

}
