package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// параметры командной строки
var (
	column        int
	number        bool
	reverse       bool
	unique        bool
	month         bool
	ignoreTBlanks bool
	checkSort     bool
	sizeNumber    bool
)

// инициализация флагов из командной строки
func init() {
	flag.IntVar(&column, "k", 0, "sort by column")
	flag.BoolVar(&number, "n", false, "sort by numeric value")
	flag.BoolVar(&reverse, "r", false, "sort in reverse order")
	flag.BoolVar(&unique, "u", false, "output only unique lines")
	flag.BoolVar(&month, "M", false, "sort by month name (Jan, Feb, etc.)")
	flag.BoolVar(&ignoreTBlanks, "b", false, "ignore trailing blanks")
	flag.BoolVar(&checkSort, "c", false, "check if data is sorted")
	flag.BoolVar(&sizeNumber, "h", false, "sort by human-readable numbers")
}

func main() {
	//считываем флаги
	flag.Parse()

	//если не указан файл в аргументах, то читаем из os.Stdin, иначе идем в файл
	var input io.Reader
	if flag.NArg() == 0 {
		input = os.Stdin
	} else {
		file, err := os.Open(flag.Arg(0))
		if err != nil {
			log.Fatalf("opening file error: %v", err)
		}
		defer file.Close()
		input = file
	}

	//читаем строки
	lines, err := readStrings(input)
	if err != nil {
		log.Fatalf("input reading error: %v", err)
	}

	// если задан флаг -c, проверяем отсортированы ли строки
	if checkSort {
		if isSorted(lines) {
			fmt.Println("File is sorted")
			os.Exit(0)
		} else {
			fmt.Println("File is not sorted")
			os.Exit(1)
		}
	}

	// выводим результат
	lines = sortStrings(lines)
	for _, line := range lines {
		fmt.Println(line)
	}
}

// чтение строк из Reader
func readStrings(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	buf := make([]byte, 0, 64*1024) // Буфер 64KB
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if ignoreTBlanks {
			line = strings.TrimRight(line, " \t")
		}
		lines = append(lines, line)
	}
	return lines, scanner.Err()
}

// проверка, отсортированы ли строки
func isSorted(lines []string) bool {
	for i := 1; i < len(lines); i++ {
		if compareStrings(lines[i-1], lines[i]) > 0 {
			return false
		}
	}
	return true
}

// сортировка с учетом флагов
func sortStrings(lines []string) []string {
	// сортировка с сохранением порядка равных элементов
	sort.SliceStable(lines, func(i, j int) bool {
		return compareStrings(lines[i], lines[j]) < 0
	})

	// удаление дубликатов, если установлен флаг -u
	if unique {
		lines = removeDuplicates(lines)
	}

	// обратный порядок, если установлен флаг -r
	if reverse {
		for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
			lines[i], lines[j] = lines[j], lines[i]
		}
	}
	return lines
}

// сравнение двух строк с учетом флагов
func compareStrings(a, b string) int {
	if column > 0 {
		a, b = getColumn(a, column), getColumn(b, column)
	}
	if month {
		return compareMonth(a, b)
	}
	if sizeNumber {
		return compareHumanReadable(a, b)
	}
	if number {
		return compareNumeric(a, b)
	}
	return strings.Compare(a, b)
}

// получение n-го столбца из строки (разделение по табуляции)
func getColumn(s string, col int) string {
	cols := strings.Split(s, "\t")
	if col <= len(cols) {
		return cols[col-1]
	}
	return ""
}

// сравнение как чисел с плавающей точкой
func compareNumeric(a, b string) int {
	af, ae := strconv.ParseFloat(a, 64)
	bf, be := strconv.ParseFloat(b, 64)
	if ae == nil && be == nil {
		if af < bf {
			return -1
		} else if af > bf {
			return 1
		}
		return 0
	}
	return strings.Compare(a, b)
}

// сравнение размеров
func compareHumanReadable(a, b string) int {
	parse := func(s string) float64 {
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			return 0
		}
		unit := s[len(s)-1]
		mult := 1.0
		switch unit {
		case 'K', 'k':
			mult = 1e3
		case 'M', 'm':
			mult = 1e6
		case 'G', 'g':
			mult = 1e9
		case 'T', 't':
			mult = 1e12
		default:
			unit = 0
		}
		if unit != 0 {
			s = s[:len(s)-1]
		}
		num, _ := strconv.ParseFloat(s, 64)
		return num * mult
	}
	af := parse(a)
	bf := parse(b)
	if af < bf {
		return -1
	} else if af > bf {
		return 1
	}
	return 0
}

func compareMonth(a, b string) int {
	months := map[string]time.Month{
		"Jan": time.January, "Feb": time.February, "Mar": time.March, "Apr": time.April,
		"May": time.May, "Jun": time.June, "Jul": time.July, "Aug": time.August,
		"Sep": time.September, "Oct": time.October, "Nov": time.November, "Dec": time.December,
	}
	//приводим к месяцу
	a = strings.Title(strings.ToLower(a[:3]))
	b = strings.Title(strings.ToLower(b[:3]))
	ma, oka := months[a]
	mb, okb := months[b]
	if oka && okb {
		if ma < mb {
			return -1
		} else if ma > mb {
			return 1
		}
		return 0
	}

	//если месяц не определен, сравниваем как строки
	return strings.Compare(a, b)
}

// Удаление дубликатов из отсортированного списка
func removeDuplicates(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}
	j := 0
	for i := 1; i < len(lines); i++ {
		if compareStrings(lines[j], lines[i]) != 0 {
			j++
			lines[j] = lines[i]
		}
	}
	return lines[:j+1]
}
