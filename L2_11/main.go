package main

import (
	"fmt"
	"sort"
	"strings"
)

func findAnagrams(words []string) map[string][]string {
	anagramMp := make(map[string][]string) // первое попавшееся слово: {анаграммы}
	keyMp := make(map[string]string)       // сортированное слово: первое попавшееся слово

	for _, word := range words {
		finalWord := strings.ToLower(word)
		chars := strings.Split(finalWord, "")
		// сортируем строку посимвольно
		sort.Strings(chars)
		strings.Join(chars, "")
		sorted := strings.Join(chars, "")

		//смотрим в keyMp: если такая последовательность символов была, то берем ключ (fst) и добавляем в anagramMp
		// если последовательности до этого не было, то записываем в keyMp и инициализируем слайс с новым словом
		if fst, ok := keyMp[sorted]; ok {
			anagramMp[fst] = append(anagramMp[fst], finalWord)
		} else {
			keyMp[sorted] = finalWord
			anagramMp[finalWord] = []string{finalWord}
		}
	}

	// Удаляем одинаковые элементы
	result := make(map[string][]string)
	for key, strArr := range anagramMp {
		if len(strArr) > 1 {
			sort.Strings(strArr)
			// Удаляем дубликаты (тут они идут подряд, так как отсортированы)
			uniqueWords := removeDuplicates(strArr)
			result[key] = uniqueWords
		}
	}
	// p.s. асимптотическая сложность цикла выше не превышает n*m*log(m), так как в худшем случае
	// при разбиении в 1 слово мы отсортируем весь список за n*log(n).

	return result
}

// удаляем дубликаты
func removeDuplicates(words []string) []string {
	if len(words) == 0 {
		return words
	}
	// уникальные элементы
	unique := []string{words[0]}
	for i := 1; i < len(words); i++ {
		// если элементы не совпали, то добавляем в итоговый слайс
		if words[i] != words[i-1] {
			unique = append(unique, words[i])
		}
	}
	return unique
}

func main() {
	//тесты
	tests := [][]string{
		{
			"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол",
		},
		{
			"Кот", "ток", "окТ", "кто",
		},
		{
			"лиса", "сила", "лиса", "сила", "лиса",
		},
		{
			"рот", "тор", "кот", "метро",
		},
	}
	//результат
	for _, test := range tests {
		fmt.Printf("Input: %v\n", test)
		result := findAnagrams(test)
		fmt.Println("Result:")
		for key, words := range result {
			fmt.Printf("  - %q: %v\n", key, words)
		}
		fmt.Println()
	}
}
