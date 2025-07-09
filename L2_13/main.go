package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// преобразует строку параметров в массив номеров выбранных колонок (полей)
// например, вот так:  fieldsStr = "1,3-5,7" -> [1, 3, 4, 5, 7]
// возвращает ошибку если строка имеет неверный формат
func parseCols(inputCols string) ([]int, error) {
	var colsNumber []int
	// делим строку на части по запятым
	parts := strings.Split(inputCols, ",")

	for _, part := range parts {
		// обрабатываем диапазоны
		if strings.Contains(part, "-") {
			// разбиваем диапазон на начальное и конечное значения
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid range: %s", part)
			}

			// находим начало диапазона
			start, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				return nil, fmt.Errorf("invalid range start: %s", rangeParts[0])
			}

			// находим конец диапазона
			end, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid range end: %s", rangeParts[1])
			}

			// проверяем корректность диапазона
			if start > end {
				return nil, fmt.Errorf("invalid range: start > end")
			}

			// добавляем все числа из диапазона в результирующий массив
			for i := start; i <= end; i++ {
				colsNumber = append(colsNumber, i)
			}
		} else {
			//обрабатываем отдельное число (не диапазон)
			field, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid field number: %s", part)
			}
			colsNumber = append(colsNumber, field)
		}
	}

	return colsNumber, nil
}

// обрабатывает одну строку по заданным параметрам
// cols - массив номеров полей для вывода (нумерация с 1)
// separated - флаг, указывающий пропускать строки без разделителя
// str - обрабатываемая строка
// checker - разделитель полей
// возвращает обработанную строку или пустую, если строка пропущена
func processStr(str string, checker string, colsNumber []int, separated bool) string {
	// если установлен флаг -s и строка не содержит разделитель - пропускаем ее
	if separated && !strings.Contains(str, checker) {
		return ""
	}

	// разбиваем строку на колонки используя параметр разделителя
	columns := strings.Split(str, checker)
	var result []string

	// проходимся по всем выбранным полям
	for _, colNum := range colsNumber {
		// индекс нумеруется с нуля
		index := colNum - 1
		// Проверяем что индекс в пределах количества колонок и добавляем в рельутат
		if index >= 0 && index < len(columns) {
			result = append(result, columns[index])
		}
	}

	// Собираем результат обратно в строку с тем же разделителем
	return strings.Join(result, checker)
}

func main() {
	// Определяем и парсим флаги командной строки
	colsFlag := flag.String("f", "", "fields to select (e.g. 1,3-5)")
	checkerFlag := flag.String("d", "\t", "delimiter character")
	separatedFlag := flag.Bool("s", false, "only output lines containing delimiter")
	flag.Parse()

	// нет смысла запускать без флага -f, поэтому останавливаем
	if *colsFlag == "" {
		fmt.Fprintln(os.Stderr, "fields flag (-f) is required")
		os.Exit(1)
	}

	//добавляем строку с полями в массив чисел
	fields, err := parseCols(*colsFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing fields: %v\n", err)
		os.Exit(1)
	}

	//читаем построчно
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// Читаем очередную строку
		str := scanner.Text()
		// Обрабатываем строку согласно параметрам
		processed := processStr(str, *checkerFlag, fields, *separatedFlag)
		// Если строка не пустая - выводим ее
		if processed != "" {
			fmt.Println(processed)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
		os.Exit(1)
	}
}
