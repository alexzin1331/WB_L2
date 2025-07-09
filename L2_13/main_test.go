package main

import (
	"slices"
	"testing"
)

func TestParseCols(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []int
		wantErr bool
	}{
		{"Single column", "2", []int{2}, false},
		{"Multiple columns", "2,4,6", []int{2, 4, 6}, false},
		{"Simple range", "2-4", []int{2, 3, 4}, false},
		{"Mixed list and range", "2,4-6", []int{2, 4, 5, 6}, false},
		{"Invalid column format", "x", nil, true},
		{"Reversed range", "5-2", nil, true},
		{"Empty string input", "", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCols(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCols() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("parseCols() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessStr(t *testing.T) {
	tests := []struct {
		name          string
		line          string
		delimiter     string
		fields        []int
		separatedOnly bool
		want          string
	}{
		{"Select second field", "x y z", " ", []int{2}, false, "y"},
		{"Tab-separated input", "x\ty\tz", "\t", []int{1}, false, "x"},
		{"Comma delimiter with multiple fields", "x,y,z", ",", []int{2, 3}, false, "y,z"},
		{"Field index out of bounds", "m n o", " ", []int{7}, false, ""},
		{"Separated-only mode with delimiter present", "red,green,blue", ",", []int{3}, true, "blue"},
		{"Separated-only mode with no delimiter", "foobar", ",", []int{1}, true, ""},
		{"Non-separated line without -s flag", "hello", ",", []int{1}, false, "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processStr(tt.line, tt.delimiter, tt.fields, tt.separatedOnly)
			if got != tt.want {
				t.Errorf("processStr() = %q, want %q", got, tt.want)
			}
		})
	}
}
