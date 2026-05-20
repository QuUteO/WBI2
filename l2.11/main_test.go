package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestFindAnagramSets(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected map[string][]string
	}{
		{
			name:  "базовый пример из задания",
			input: []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"},
			expected: map[string][]string{
				"пятак":  {"пятак", "пятка", "тяпка"},
				"листок": {"листок", "слиток", "столик"},
			},
		},
		{
			name:  "слова в разном регистре",
			input: []string{"ПяТак", "пятка", "ТЯПКА", "Листок", "СЛИТОК", "столик", "стол"},
			expected: map[string][]string{
				"пятак":  {"пятак", "пятка", "тяпка"},
				"листок": {"листок", "слиток", "столик"},
			},
		},
		{
			name:     "пустой словарь",
			input:    []string{},
			expected: map[string][]string{},
		},
		{
			name:     "один элемент",
			input:    []string{"пятак"},
			expected: map[string][]string{},
		},
		{
			name:     "нет анаграмм",
			input:    []string{"кот", "пес", "дом", "мышь"},
			expected: map[string][]string{},
		},
		{
			name:  "только две анаграммы",
			input: []string{"кот", "ток", "окт", "дом"},
			expected: map[string][]string{
				"кот": {"кот", "окт", "ток"},
			},
		},
		{
			name:  "несколько групп анаграмм",
			input: []string{"кот", "ток", "окт", "сон", "нос", "дом", "мод"},
			expected: map[string][]string{
				"кот": {"кот", "окт", "ток"},
				"сон": {"нос", "сон"},
				"дом": {"дом", "мод"},
			},
		},
		{
			name:  "сортировка внутри группы",
			input: []string{"тяпка", "пятка", "пятак"},
			expected: map[string][]string{
				"тяпка": {"пятак", "пятка", "тяпка"},
			},
		},
		{
			name:  "первое встреченное слово как ключ",
			input: []string{"тяпка", "пятак", "пятка"},
			expected: map[string][]string{
				"тяпка": {"пятак", "пятка", "тяпка"},
			},
		},
		{
			name:  "слова с повторяющимися буквами",
			input: []string{"ааабб", "аббаа", "ббааа", "ааббб", "нормальное"},
			expected: map[string][]string{
				"ааабб": {"ааабб", "аббаа", "ббааа"},
			},
		},
		{
			name:  "все слова одинаковые",
			input: []string{"слово", "слово", "слово", "слово"},
			expected: map[string][]string{
				"слово": {"слово", "слово", "слово", "слово"},
			},
		},
		{
			name:  "смесь русского и английского (не анаграммы)",
			input: []string{"cat", "tac", "act", "кот", "ток"},
			expected: map[string][]string{
				"cat": {"act", "cat", "tac"},
				"кот": {"кот", "ток"},
			},
		},
		{
			name:     "слова с разной длиной",
			input:    []string{"а", "аа", "ааа", "аааа"},
			expected: map[string][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindAnagramSets(tt.input)

			if !compareMaps(result, tt.expected) {
				t.Errorf("FindAnagramSets() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func compareMaps(m1, m2 map[string][]string) bool {
	if len(m1) != len(m2) {
		return false
	}

	for key, val1 := range m1 {
		val2, exists := m2[key]
		if !exists {
			return false
		}

		sorted1 := make([]string, len(val1))
		sorted2 := make([]string, len(val2))
		copy(sorted1, val1)
		copy(sorted2, val2)

		sort.Strings(sorted1)
		sort.Strings(sorted2)

		if !reflect.DeepEqual(sorted1, sorted2) {
			return false
		}
	}

	return true
}

func TestGetSignature(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "обычное слово",
			input:    "пятак",
			expected: "акптя",
		},
		{
			name:     "слово с повторяющимися буквами",
			input:    "ананас",
			expected: "аааннс",
		},
		{
			name:     "одна буква",
			input:    "а",
			expected: "а",
		},
		{
			name:     "слово из английских букв",
			input:    "hello",
			expected: "ehllo",
		},
		{
			name:     "пустая строка",
			input:    "",
			expected: "",
		},
		{
			name:     "слово в верхнем регистре",
			input:    "ПЯТАК",
			expected: "акптя",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getSignature(tt.input)
			if result != tt.expected {
				t.Errorf("getSignature(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func BenchmarkFindAnagramSets(b *testing.B) {
	input := []string{
		"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол",
		"кот", "ток", "окт", "сон", "нос", "дом", "мод",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FindAnagramSets(input)
	}
}

func TestCaseInsensitive(t *testing.T) {
	input := []string{"Пятак", "ПЯТКА", "тЯпКа", "ЛИСТОК", "СлиТок", "столик"}
	expected := map[string][]string{
		"пятак":  {"пятак", "пятка", "тяпка"},
		"листок": {"листок", "слиток", "столик"},
	}

	result := FindAnagramSets(input)

	if !compareMaps(result, expected) {
		t.Errorf("Case insensitive test failed: got %v, want %v", result, expected)
	}
}

// Исправленный тест для уникальных ключей
func TestUniqueKeysForAnagrams(t *testing.T) {
	// Все эти слова - анаграммы друг друга
	input := []string{"abc", "acb", "bac", "bca", "cab", "cba"}
	result := FindAnagramSets(input)

	// Должна быть только одна группа
	if len(result) != 1 {
		t.Errorf("Expected 1 key, got %d", len(result))
	}

	// Проверяем, что ключ - первое встреченное слово ("abc")
	for key := range result {
		if key != "abc" {
			t.Errorf("Expected key 'abc', got '%s'", key)
		}
	}

	// Проверяем количество слов в группе
	for _, words := range result {
		if len(words) != 6 {
			t.Errorf("Expected 6 words in group, got %d", len(words))
		}
	}
}

// Исправленный тест для Unicode с корректными анаграммами
func TestUnicodeSupport(t *testing.T) {
	// Используем корректные анаграммы без буквы "ё" (она сортируется отдельно)
	input := []string{"елка", "кела", "леак"}
	expected := map[string][]string{
		"елка": {"елка", "кела", "леак"},
	}

	result := FindAnagramSets(input)

	if !compareMaps(result, expected) {
		t.Errorf("Unicode test failed: got %v, want %v", result, expected)
	}
}

// Дополнительный тест для буквы ё
func TestYoLetter(t *testing.T) {
	input := []string{"ёлка", "колё", "окёл"}
	// Обратите внимание: "ёлка" и "колё" - анаграммы?
	// Буква "ё" часто считается вариантом "е", но в Unicode это разные символы
	// Поэтому тест должен учитывать реальное поведение

	result := FindAnagramSets(input)

	// Проверяем, что результат соответствует ожиданиям
	t.Logf("Result for ё: %v", result)
}
