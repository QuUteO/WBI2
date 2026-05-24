package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// FieldSet представляет множество номеров полей для вывода
type FieldSet map[int]bool

// parseFieldSpec парсит спецификацию полей вида "1,3-5,7"
func parseFieldSpec(spec string) (FieldSet, error) {
	fields := make(FieldSet)

	if spec == "" {
		return nil, fmt.Errorf("не указано ни одного поля")
	}

	parts := strings.Split(spec, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, "-") {
			// Диапазон: start-end
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("некорректный диапазон полей: %s", part)
			}

			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return nil, fmt.Errorf("некорректный номер поля: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return nil, fmt.Errorf("некорректный номер поля: %s", rangeParts[1])
			}

			if start < 1 {
				return nil, fmt.Errorf("номер поля должен быть >= 1, получено %d", start)
			}
			if end < 1 {
				return nil, fmt.Errorf("номер поля должен быть >= 1, получено %d", end)
			}

			if start > end {
				start, end = end, start
			}

			for i := start; i <= end; i++ {
				fields[i] = true
			}
		} else {
			// Одиночное поле
			num, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("некорректный номер поля: %s", part)
			}
			if num < 1 {
				return nil, fmt.Errorf("номер поля должен быть >= 1, получено %d", num)
			}
			fields[num] = true
		}
	}

	if len(fields) == 0 {
		return nil, fmt.Errorf("не указано ни одного поля")
	}

	return fields, nil
}

// getSortedFields возвращает отсортированный список номеров полей
func getSortedFields(fields FieldSet) []int {
	result := make([]int, 0, len(fields))
	for f := range fields {
		result = append(result, f)
	}
	sort.Ints(result)
	return result
}

// processLine обрабатывает одну строку и возвращает результат
func processLine(line string, fields FieldSet, delimiter string, separatedOnly bool) (string, bool) {
	// Проверяем, содержит ли строка разделитель
	if separatedOnly && !strings.Contains(line, delimiter) {
		return "", false
	}

	// Разделяем строку
	parts := strings.Split(line, delimiter)

	// Собираем нужные поля
	var selected []string
	sortedFields := getSortedFields(fields)

	for _, fieldNum := range sortedFields {
		// Номера полей в спецификации 1-индексированные
		idx := fieldNum - 1
		if idx < len(parts) {
			selected = append(selected, parts[idx])
		}
		// Поля за пределами границ игнорируем
	}

	if len(selected) == 0 {
		return "", false
	}

	return strings.Join(selected, delimiter), true
}

func main() {
	// Парсинг аргументов командной строки
	var (
		fieldsSpec    string
		delimiter     string
		separatedOnly bool
	)

	flag.StringVar(&fieldsSpec, "f", "", "номера полей для вывода (например: 1,3-5,7)")
	flag.StringVar(&fieldsSpec, "fields", "", "номера полей для вывода (например: 1,3-5,7)")
	flag.StringVar(&delimiter, "d", "\t", "разделитель полей (по умолчанию: табуляция)")
	flag.StringVar(&delimiter, "delimiter", "\t", "разделитель полей (по умолчанию: табуляция)")
	flag.BoolVar(&separatedOnly, "s", false, "выводить только строки, содержащие разделитель")
	flag.BoolVar(&separatedOnly, "separated-only", false, "выводить только строки, содержащие разделитель")

	flag.Parse()

	// Проверка обязательного параметра -f
	if fieldsSpec == "" {
		fmt.Fprintf(os.Stderr, "Ошибка: требуется указать параметр -f/--fields\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Парсинг спецификации полей
	fields, err := parseFieldSpec(fieldsSpec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка в параметре -f: %v\n", err)
		os.Exit(1)
	}

	// Чтение из STDIN
	scanner := bufio.NewScanner(os.Stdin)
	// Увеличиваем буфер для больших строк
	const maxCapacity = 1024 * 1024 // 1 MB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		result, ok := processLine(line, fields, delimiter, separatedOnly)
		if ok {
			fmt.Println(result)
			lineCount++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка при чтении ввода: %v\n", err)
		os.Exit(1)
	}
}
