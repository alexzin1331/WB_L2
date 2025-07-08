package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

type config struct {
	after      int
	before     int
	around     int
	count      bool
	ignoreCase bool
	invert     bool
	fixed      bool
	lineNum    bool
	pattern    string
	filename   string
}

func main() {
	cfg := parseFlags()

	var input io.Reader
	if cfg.filename == "" || cfg.filename == "-" {
		input = os.Stdin
	} else {
		file, err := os.Open(cfg.filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "grep: %s: %v\n", cfg.filename, err)
			os.Exit(1)
		}
		defer file.Close()
		input = file
	}

	err := grep(input, os.Stdout, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "grep: %v\n", err)
		os.Exit(1)
	}
}

func parseFlags() config {
	var cfg config

	flag.IntVar(&cfg.after, "A", 0, "after each found line, additionally output N lines after it")
	flag.IntVar(&cfg.before, "B", 0, "output N lines before each found line")
	flag.IntVar(&cfg.around, "C", 0, "output N lines of context around the found string")
	flag.BoolVar(&cfg.count, "c", false, "output only the number of lines that match the template")
	flag.BoolVar(&cfg.ignoreCase, "i", false, "ignore case distinctions")
	flag.BoolVar(&cfg.invert, "v", false, " invert the filter: output lines that do not contain a template")
	flag.BoolVar(&cfg.fixed, "F", false, "treat the pattern as a fixed string rather than a regular expression")
	flag.BoolVar(&cfg.lineNum, "n", false, "print the line number before each found line")

	// парсим флаги
	flag.Parse()

	// если задан -C, он переопределяет -A и -B
	if cfg.around > 0 {
		cfg.after = cfg.around
		cfg.before = cfg.around
	}

	args := flag.Args()
	if len(args) < 1 {
		log.Fatalln("not enough arguments")
	}

	// первый аргумент является шаблоном
	cfg.pattern = args[0]
	if len(args) > 1 {
		// второй аргумент это файл, если есть
		cfg.filename = args[1]
	}

	return cfg
}

func grep(input io.Reader, output io.Writer, cfg config) error {
	scanner := bufio.NewScanner(input)
	var pattern *regexp.Regexp
	var err error

	// выполняем точное выражение строки
	if cfg.fixed {
		// игнорируем регулярные выражения (добавляем \)
		patternString := regexp.QuoteMeta(cfg.pattern)
		if cfg.ignoreCase {
			pattern, err = regexp.Compile("(?i)" + patternString) // Игнорируем регистр
		} else {
			pattern, err = regexp.Compile(patternString)
		}
		//учитываем регулярные выражения
	} else {
		if cfg.ignoreCase {
			pattern, err = regexp.Compile("(?i)" + cfg.pattern) // Игнорируем регистр
		} else {
			pattern, err = regexp.Compile(cfg.pattern)
		}
	}
	if err != nil {
		return fmt.Errorf("invalid pattern: %v", err)
	}

	// Сохраняем все строки и их номера
	inputStr := make([]string, 0)
	strNumbers := make([]int, 0)
	strIdx := 1
	for scanner.Scan() {
		inputStr = append(inputStr, scanner.Text())
		strNumbers = append(strNumbers, strIdx)
		strIdx++
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	//ищем необходимые строки, соответствующие флагу с учетом флага -v
	var suit []int
	for i, line := range inputStr {
		matched := pattern.MatchString(line)
		if cfg.invert {
			matched = !matched
		}
		if matched {
			suit = append(suit, i)
		}
	}

	//если требуется вывести количество совпавших строк
	if cfg.count {
		fmt.Fprintf(output, "%d\n", len(suit))
		return nil
	}

	checkPrint := make(map[int]struct{})
	last := -1

	// печатаем результат с учётом флагов (-A, -B, -C)
	for _, match := range suit {
		start := max(0, match-cfg.before)            // начало интервала
		end := min(len(inputStr)-1, match+cfg.after) // конец интервала

		// если контексты пересекаются, корректируем start
		if last >= 0 && start <= last {
			start = last + 1
		}

		// Печатаем строки в интервале
		for i := start; i <= end; i++ {
			//проверяем, была ли напечатана строка
			if _, ok := checkPrint[i]; ok {
				continue
			}
			checkPrint[i] = struct{}{}

			//Если флаг -n активен, добавляем к строке её номер в формате "номер:"
			var prefix string
			if cfg.lineNum {
				prefix = fmt.Sprintf("%d:", strNumbers[i])
			}

			line := inputStr[i]
			if i == match { // подсвечиваем совпадение
				line = pattern.ReplaceAllString(line, ">>$0<<")
			}
			fmt.Fprintf(output, "%s%s\n", prefix, line)
		}
		last = end
	}

	return nil
}
