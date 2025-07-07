package main

import (
	"reflect"
	"testing"
)

func TestSortFlags(t *testing.T) {
	tests := []struct {
		name   string
		input  []string
		flags  func()
		expect []string
	}{
		{
			//обычные числа
			name:  "-n (numeric sort)",
			input: []string{"10", "2", "33"},
			flags: func() {
				number, reverse, unique, month, sizeNumber, column = true, false, false, false, false, 0
			},
			expect: []string{"2", "10", "33"},
		},
		{
			//обратная сортировка
			name:  "-r (reverse sort)",
			input: []string{"apple", "banana", "cherry"},
			flags: func() {
				number, reverse, unique, month, sizeNumber, column = false, true, false, false, false, 0
			},
			expect: []string{"cherry", "banana", "apple"},
		},
		{
			//очистка от дубликатов
			name:  "-u (unique)",
			input: []string{"a", "b", "a", "c"},
			flags: func() {
				number, reverse, unique, month, sizeNumber, column = false, false, true, false, false, 0
			},
			expect: []string{"a", "b", "c"},
		},
		{
			//сортировка по дате
			name:  "-M (month sort)",
			input: []string{"Feb", "Jan", "Dec"},
			flags: func() {
				number, reverse, unique, month, sizeNumber, column = false, false, false, true, false, 0
			},
			expect: []string{"Jan", "Feb", "Dec"},
		},
		{
			// Смешанные единицы измерения - должны сортироваться
			// по числовому эквиваленту
			name:  "-h (human-readable sort)",
			input: []string{"1K", "200", "3M"},
			flags: func() {
				number, reverse, unique, month, sizeNumber, column = false, false, false, false, true, 0
			},
			expect: []string{"200", "1K", "3M"},
		},
		{
			// Отсутствующие колонки - строки без достаточного
			// количества колонок должны считаться пустыми
			name:  "-k2 (sort by 2nd column)",
			input: []string{"a	3", "b	1", "c	2"},
			flags: func() {
				number, reverse, unique, month, sizeNumber, column = false, false, false, false, false, 2
			},
			expect: []string{"b\t1", "c\t2", "a\t3"},
		},
		{
			// Пустые числовые поля - должны обрабатываться
			// как нули или пустые значения
			name:  "-k2nr (numeric reverse by column)",
			input: []string{"x	10", "y	2", "z	5"},
			flags: func() {
				number, reverse, unique, month, sizeNumber, column = true, true, false, false, false, 2
			},
			expect: []string{"x\t10", "z\t5", "y\t2"},
		},
		{
			// Комбинация флагов - проверка комплексного поведения
			// числовая обратная сортировка по колонке с удалением дублей
			name:  "-k2nu (unique numeric by column)",
			input: []string{"a	1", "b	2", "c	1", "d	3"},
			flags: func() {
				number, reverse, unique, month, sizeNumber, column = true, false, true, false, false, 2
			},
			expect: []string{"a\t1", "b\t2", "d\t3"},
		},
		{
			//игнорируем хвостовые пробелы
			name:  "-b with trailing spaces",
			input: []string{"a  ", " b", "  c"},
			flags: func() {
				ignoreTBlanks = true
				number, reverse, unique, month, sizeNumber, column = false, false, false, false, false, 0
			},
			expect: []string{"  c", " b", "a  "},
		},
		{
			//проверка должна быть в main
			name:  "-c with sorted input",
			input: []string{"a", "b", "c"},
			flags: func() {
				checkSort = true
			},
			expect: []string{"a", "b", "c"}, // Actual check happens in main()
		},
		{
			//проверка без параметров
			name:  "single line input",
			input: []string{"single line"},
			flags: func() {
				number, reverse, unique, month, sizeNumber, column = false, false, false, false, false, 0 //все false
			},
			expect: []string{"single line"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.flags()
			lines := append([]string(nil), tt.input...)
			lines = sortStrings(lines)
			if !reflect.DeepEqual(lines, tt.expect) {
				t.Errorf("got %v, want %v", lines, tt.expect)
			}
		})
	}
}
