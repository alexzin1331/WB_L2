package main

import (
	"testing"
)

func TestResolveString(t *testing.T) {
	tests := []struct {
		number   int
		name     string
		input    string
		expected string
		hasError bool
	}{
		// из условия задачи
		{number: 1, name: "a4bc2d5e", input: "a4bc2d5e", expected: "aaaabccddddde", hasError: false},
		{number: 2, name: "abcd", input: "abcd", expected: "abcd", hasError: false},
		{number: 3, name: "45", input: "45", expected: "", hasError: true},
		{number: 4, name: "empty", input: "", expected: "", hasError: false},

		// c Escape
		{number: 5, name: `qwe\4\5`, input: `qwe\4\5`, expected: "qwe45", hasError: false},
		{number: 6, name: `qwe\45`, input: `qwe\45`, expected: "qwe44444", hasError: false},
		{number: 7, name: `abc\`, input: `abc\`, expected: "", hasError: true}, // Ошибка: строка заканчивается на \

		// дополнительные тесты
		{number: 8, name: "a", input: "a", expected: "a", hasError: false},
		{number: 9, name: "a10b", input: "a10b", expected: "aaaaaaaaaab", hasError: false},
		{number: 10, name: "0abc", input: "0abc", expected: "", hasError: true}, // Ошибка: начинается с цифры
		{number: 11, name: `\4`, input: `\4`, expected: "4", hasError: false},

		// русские буквы
		{number: 12, name: "ярусский", input: `ярусский`, expected: "ярусский", hasError: false},
		{number: 13, name: "л2к5", input: `л2к5`, expected: "ллккккк", hasError: false},
		{number: 14, name: `\4`, input: `\4`, expected: `4`, hasError: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := resolveString(tt.input)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error, but got none (input: %q)", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v (input: %q)", err, tt.input)
				} else if res != tt.expected {
					t.Errorf("Expected %q, got %q (input: %q)", tt.expected, res, tt.input)
				}
			}
		})
	}
}
