package main

import (
	"strings"
	"testing"
)

// TestSortBasic проверяет базовую сортировку
func TestSortBasic(t *testing.T) {
	lines := []string{"banana", "apple", "cherry"}
	config := &ConfigSort{}

	result := sortLines(lines, config)

	expected := []string{"apple", "banana", "cherry"}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("Expected %q, got %q", expected[i], result[i])
		}
	}
}

// TestSortReverse проверяет флаг -r
func TestSortReverse(t *testing.T) {
	lines := []string{"apple", "banana", "cherry"}
	config := &ConfigSort{reverse: true}

	result := sortLines(lines, config)

	expected := []string{"cherry", "banana", "apple"}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("Expected %q, got %q", expected[i], result[i])
		}
	}
}

// TestSortNumeric проверяет флаг -n
func TestSortNumeric(t *testing.T) {
	lines := []string{"10", "2", "30", "1", "20"}
	config := &ConfigSort{numeric: true}

	result := sortLines(lines, config)

	expected := []string{"1", "2", "10", "20", "30"}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("Expected %q, got %q", expected[i], result[i])
		}
	}
}

// TestSortNumericWithText проверяет числовую сортировку со смешанными строками
func TestSortNumericWithText(t *testing.T) {
	lines := []string{"file10", "file2", "file30", "file1"}
	config := &ConfigSort{numeric: true}

	result := sortLines(lines, config)

	expected := []string{"file1", "file2", "file10", "file30"}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("Expected %q, got %q", expected[i], result[i])
		}
	}
}

// TestSortUnique проверяет флаг -u
func TestSortUnique(t *testing.T) {
	lines := []string{"apple", "banana", "apple", "cherry", "banana"}
	config := &ConfigSort{unique: true}

	result := sortLines(lines, config)

	expected := []string{"apple", "banana", "cherry"}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("Expected %q, got %q", expected[i], result[i])
		}
	}
}

// TestSortByColumn проверяет флаг -k
func TestSortByColumn(t *testing.T) {
	lines := []string{
		"a\t3\tx",
		"b\t1\ty",
		"c\t2\tz",
	}
	config := &ConfigSort{column: 2, numeric: true}

	result := sortLines(lines, config)

	expected := []string{
		"b\t1\ty",
		"c\t2\tz",
		"a\t3\tx",
	}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("Expected %q, got %q", expected[i], result[i])
		}
	}
}

// TestSortByColumnText проверяет сортировку по колонке как текст
func TestSortByColumnText(t *testing.T) {
	lines := []string{
		"a\tcherry\tx",
		"b\tapple\ty",
		"c\tbanana\tz",
	}
	config := &ConfigSort{column: 2}

	result := sortLines(lines, config)

	expected := []string{
		"b\tapple\ty",
		"c\tbanana\tz",
		"a\tcherry\tx",
	}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("Expected %q, got %q", expected[i], result[i])
		}
	}
}

// TestSortIgnoreBlanks проверяет флаг -b
func TestSortIgnoreBlanks(t *testing.T) {
	lines := []string{"banana   ", "apple", "cherry   "}
	config := &ConfigSort{ignoreBlanks: true}

	result := sortLines(lines, config)

	// Строки должны сравниваться без учёта пробелов
	expected := []string{"apple", "banana   ", "cherry   "}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("Expected %q, got %q", expected[i], result[i])
		}
	}
}

// TestSortIgnoreBlanksWithUnique проверяет -b и -u вместе
func TestSortIgnoreBlanksWithUnique(t *testing.T) {
	lines := []string{"apple   ", "banana", "apple", "cherry"}
	config := &ConfigSort{ignoreBlanks: true, unique: true}

	result := sortLines(lines, config)

	expected := []string{"apple   ", "banana", "cherry"}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("Expected %q, got %q", expected[i], result[i])
		}
	}
}

// TestSortMonth проверяет флаг -M
func TestSortMonth(t *testing.T) {
	lines := []string{"Mar", "Jan", "Feb", "Apr", "Dec"}
	config := &ConfigSort{month: true}

	result := sortLines(lines, config)

	expected := []string{"Jan", "Feb", "Mar", "Apr", "Dec"}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("Expected %q, got %q", expected[i], result[i])
		}
	}
}

// TestSortMonthFullNames проверяет полные названия месяцев
func TestSortMonthFullNames(t *testing.T) {
	lines := []string{"March", "January", "February", "April"}
	config := &ConfigSort{month: true}

	result := sortLines(lines, config)

	expected := []string{"January", "February", "March", "April"}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("Expected %q, got %q", expected[i], result[i])
		}
	}
}

// TestSortHuman проверяет флаг -h
func TestSortHuman(t *testing.T) {
	lines := []string{"1M", "500K", "2G", "1K"}
	config := &ConfigSort{human: true}

	result := sortLines(lines, config)

	expected := []string{"1K", "500K", "1M", "2G"}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("Expected %q, got %q", expected[i], result[i])
		}
	}
}

// TestSortHumanWithDecimals проверяет human с десятичными числами
func TestSortHumanWithDecimals(t *testing.T) {
	lines := []string{"1.5K", "2M", "0.5K", "1GB"}
	config := &ConfigSort{human: true}

	result := sortLines(lines, config)

	expected := []string{"0.5K", "1.5K", "2M", "1GB"}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("Expected %q, got %q", expected[i], result[i])
		}
	}
}

// TestSortCombined проверяет комбинацию флагов
func TestSortCombined(t *testing.T) {
	lines := []string{
		"a\t10\tx",
		"b\t2\ty",
		"c\t30\tz",
		"d\t2\tw",
	}
	config := &ConfigSort{column: 2, numeric: true, reverse: true, unique: true}

	result := sortLines(lines, config)

	expected := []string{
		"c\t30\tz",
		"a\t10\tx",
		"b\t2\ty",
	}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("Expected %q, got %q", expected[i], result[i])
		}
	}
}

// TestIsSorted проверяет функцию isSorted
func TestIsSorted(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		config   *ConfigSort
		expected bool
	}{
		{
			name:     "пустой слайс",
			lines:    []string{},
			config:   &ConfigSort{},
			expected: true,
		},
		{
			name:     "одна строка",
			lines:    []string{"apple"},
			config:   &ConfigSort{},
			expected: true,
		},
		{
			name:     "отсортировано",
			lines:    []string{"apple", "banana", "cherry"},
			config:   &ConfigSort{},
			expected: true,
		},
		{
			name:     "не отсортировано",
			lines:    []string{"banana", "apple", "cherry"},
			config:   &ConfigSort{},
			expected: false,
		},
		{
			name:     "числовая отсортировано",
			lines:    []string{"1", "2", "10", "20"},
			config:   &ConfigSort{numeric: true},
			expected: true,
		},
		{
			name:     "числовая не отсортировано",
			lines:    []string{"1", "10", "2", "20"},
			config:   &ConfigSort{numeric: true},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSorted(tt.lines, tt.config)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestExtractNumber проверяет извлечение чисел
func TestExtractNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"123", 123},
		{"  -45.67  ", -45.67},
		{"abc123xyz", 123},
		{"no numbers", 0},
		{"10 items", 10},
		{"-5.5", -5.5},
		{"0", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := extractNumber(tt.input)
			if result != tt.expected {
				t.Errorf("extractNumber(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestParseHumanSize проверяет парсинг human-размеров
func TestParseHumanSize(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"1K", 1024},
		{"1M", 1024 * 1024},
		{"1G", 1024 * 1024 * 1024},
		{"1.5K", 1536},
		{"500K", 500 * 1024},
		{"2MB", 2 * 1024 * 1024},
		{"1GB", 1 * 1024 * 1024 * 1024},
		{"100", 100},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseHumanSize(tt.input)
			if result != tt.expected {
				t.Errorf("parseHumanSize(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestGetColumn проверяет извлечение колонки
func TestGetColumn(t *testing.T) {
	line := "a\tb\tc"

	tests := []struct {
		column   int
		expected string
	}{
		{1, "a"},
		{2, "b"},
		{3, "c"},
		{4, ""},
	}

	for _, tt := range tests {
		t.Run("column "+string(rune(tt.column)), func(t *testing.T) {
			result := getColumn(line, tt.column)
			if result != tt.expected {
				t.Errorf("getColumn(%q, %d) = %q, want %q", line, tt.column, result, tt.expected)
			}
		})
	}
}

// TestTrimBlanks проверяет обрезку пробелов
func TestTrimBlanks(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello   ", "hello"},
		{"hello", "hello"},
		{"  hello  ", "  hello"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := trimBlanks(tt.input)
			if result != tt.expected {
				t.Errorf("trimBlanks(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestExtractMonthName проверяет извлечение названия месяца
func TestExtractMonthName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"January", "Jan"},
		{"Jan", "Jan"},
		{"February", "Feb"},
		{"Mar", "Mar"},
		{"   December   ", "Dec"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := extractMonthName(tt.input)
			if result != tt.expected {
				t.Errorf("extractMonthName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestReadLines проверяет чтение строк
func TestReadLines(t *testing.T) {
	input := "apple\nbanana\ncherry\n"
	reader := strings.NewReader(input)
	config := &ConfigSort{}

	lines, err := readLines(reader, config)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := []string{"apple", "banana", "cherry"}
	if len(lines) != len(expected) {
		t.Errorf("Expected %d lines, got %d", len(expected), len(lines))
	}
	for i := range expected {
		if lines[i] != expected[i] {
			t.Errorf("Expected %q, got %q", expected[i], lines[i])
		}
	}
}

// TestReadLinesWithBlanks проверяет чтение с игнорированием пробелов
func TestReadLinesWithBlanks(t *testing.T) {
	input := "apple   \nbanana\ncherry   \n"
	reader := strings.NewReader(input)
	config := &ConfigSort{ignoreBlanks: true}

	lines, err := readLines(reader, config)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := []string{"apple", "banana", "cherry"}
	for i := range expected {
		if lines[i] != expected[i] {
			t.Errorf("Expected %q, got %q", expected[i], lines[i])
		}
	}
}

// BenchmarkSort проверяет производительность на больших данных
func BenchmarkSort(b *testing.B) {
	// Создаём тестовые данные
	lines := make([]string, 10000)
	for i := 0; i < 10000; i++ {
		lines[i] = string(rune(10000 - i))
	}
	config := &ConfigSort{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sortLines(lines, config)
	}
}

// BenchmarkSortNumeric проверяет производительность числовой сортировки
func BenchmarkSortNumeric(b *testing.B) {
	lines := make([]string, 10000)
	for i := 0; i < 10000; i++ {
		lines[i] = string(rune(10000 - i))
	}
	config := &ConfigSort{numeric: true}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sortLines(lines, config)
	}
}
