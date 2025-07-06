package main

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

func resolveString(s string) (string, error) {
	escape := false
	var prevSymbol rune
	answer := strings.Builder{}
	runes := []rune(s) // для корректной обработки Unicode

	for i := 0; i < len(runes); i++ {
		v := runes[i]
		switch {
		case escape:
			// после \ добавляем любой символ
			answer.WriteRune(v)
			prevSymbol = v
			escape = false

		case v == '\\':
			// найден \, тогда ставим escape = true
			escape = true

		case unicode.IsDigit(v):
			if i == 0 {
				// если первый символ является цифрой, то возвращаем ошибку
				return "", errors.New("invalid string: starts with digit")
			}

			// получаем число в цикле, учитывая, что число может состоять из 2 и более разрядов
			numStr := string(v)
			j := i + 1
			for j < len(runes) && unicode.IsDigit(runes[j]) {
				numStr += string(runes[j])
				j++
			}

			// после считывания преобразуем в число
			count, err := strconv.Atoi(numStr)
			if err != nil {
				return "", errors.New("invalid number format")
			}

			if count == 0 {
				// если число 0, удаляем предыдущий символ
				if answer.Len() > 0 {
					temp := answer.String()
					answer.Reset()
					answer.WriteString(temp[:len(temp)-1])
				}
			} else if prevSymbol != 0 {
				// дублируем предыдущий символ count-1 раз, учитывая, что он уже есть 1 раз в строке
				answer.WriteString(strings.Repeat(string(prevSymbol), count-1))
			}

			// пропускаем обработанные цифры
			i = j - 1

		default:
			// обычный символ добавляем в любом случае
			answer.WriteRune(v)
			prevSymbol = v
		}
	}

	if escape {
		// строка закончилась на \, это ошибка
		return "", errors.New("invalid string: ends with escape character")
	}

	return answer.String(), nil
}
