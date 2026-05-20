package main

import (
	"fmt"
	"sort"
	"strings"
)

type group struct {
	first string
	words []string
}

func FindAnagramSets(words []string) map[string][]string {
	// Промежуточная структура: сигнатура -> {первое слово, список слов}
	groups := make(map[string]*group)

	// Для отслеживания порядка первых вхождений
	order := []string{}

	for _, word := range words {
		lowerWord := strings.ToLower(word)

		// Получаем сигнатуру (отсортированные буквы)
		signature := getSignature(lowerWord)

		// Если группа ещё не существует
		if _, exists := groups[signature]; !exists {
			// Запоминаем первое встреченное слово
			groups[signature] = &group{
				first: lowerWord,
				words: []string{lowerWord},
			}
			order = append(order, signature)
		} else {
			// Добавляем слово в существующую группу
			groups[signature].words = append(groups[signature].words, lowerWord)
		}
	}

	// Формируем результат
	result := make(map[string][]string)

	for _, sig := range order {
		g := groups[sig]
		// Пропускаем группы из одного слова
		if len(g.words) < 2 {
			continue
		}

		// Сортируем слова в группе по возрастанию
		sort.Strings(g.words)

		// Ключ - первое встреченное слово
		result[g.first] = g.words
	}

	return result
}

// getSignature возвращает отсортированную строку из букв слова (в нижнем регистре)
func getSignature(word string) string {
	// Приводим к нижнему регистру
	lowerWord := strings.ToLower(word)

	// Преобразуем строку в срез рун (для корректной работы с Unicode)
	runes := []rune(lowerWord)

	// Сортируем руны
	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})

	return string(runes)
}

func main() {
	input := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}

	result := FindAnagramSets(input)

	// Вывод результата
	for key, words := range result {
		fmt.Printf("%q: %q\n", key, words)
	}
}
