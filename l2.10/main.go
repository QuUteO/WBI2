package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// ConfigSort содержит все параметры сортировки
type ConfigSort struct {
	column       int
	numeric      bool
	reverse      bool
	unique       bool
	month        bool
	ignoreBlanks bool
	check        bool
	human        bool
}

// monthMap для сортировки по названиям месяцев
var monthMap = map[string]int{
	"jan": 1, "feb": 2, "mar": 3, "apr": 4, "may": 5, "jun": 6,
	"jul": 7, "aug": 8, "sep": 9, "oct": 10, "nov": 11, "dec": 12,
}

func main() {
	config := parseFlag()

	lines, err := readInputData(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sort: %v\n", err)
		os.Exit(1)
	}

	// Режим проверки (-c)
	if config.check {
		if !isSorted(lines, config) {
			fmt.Fprintln(os.Stderr, "sort: lines are not sorted")
			os.Exit(1)
		}
		return
	}

	// Сортировка
	sorted := sortLines(lines, config)

	// Вывод результата
	for _, line := range sorted {
		fmt.Println(line)
	}
}

// parseFlag парсит флаги командной строки
func parseFlag() *ConfigSort {
	var config ConfigSort

	flag.IntVar(&config.column, "k", 0, "сортировать по столбцу N")
	flag.BoolVar(&config.numeric, "n", false, "сортировать по числовому значению")
	flag.BoolVar(&config.reverse, "r", false, "сортировать в обратном порядке")
	flag.BoolVar(&config.unique, "u", false, "не выводить повторяющиеся строки")
	flag.BoolVar(&config.month, "M", false, "сортировать по названию месяца")
	flag.BoolVar(&config.ignoreBlanks, "b", false, "игнорировать хвостовые пробелы")
	flag.BoolVar(&config.check, "c", false, "проверить, отсортированы ли данные")
	flag.BoolVar(&config.human, "h", false, "сортировать с учётом суффиксов (K, M, G)")

	flag.Parse()

	return &config
}

// sortLines сортирует строки согласно конфигурации
func sortLines(lines []string, config *ConfigSort) []string {
	result := make([]string, len(lines))
	copy(result, lines)

	// Выбираем тип сортировки
	if config.month {
		// Сортировка по месяцам
		sort.Slice(result, func(i, j int) bool {
			valI := result[i]
			valJ := result[j]

			if config.ignoreBlanks {
				valI = trimBlanks(valI)
				valJ = trimBlanks(valJ)
			}

			if config.column > 0 {
				valI = getColumn(valI, config.column)
				valJ = getColumn(valJ, config.column)
			}

			monthI := extractMonthName(valI)
			monthJ := extractMonthName(valJ)

			numI := monthMap[strings.ToLower(monthI)]
			numJ := monthMap[strings.ToLower(monthJ)]

			// Если не месяц, используем обычное сравнение
			if numI == 0 && numJ == 0 {
				if config.reverse {
					return valI > valJ
				}
				return valI < valJ
			}

			if config.reverse {
				return numI > numJ
			}
			return numI < numJ
		})
	} else if config.numeric {
		// Числовая сортировка
		sort.Slice(result, func(i, j int) bool {
			valI := result[i]
			valJ := result[j]

			if config.ignoreBlanks {
				valI = trimBlanks(valI)
				valJ = trimBlanks(valJ)
			}

			if config.column > 0 {
				valI = getColumn(valI, config.column)
				valJ = getColumn(valJ, config.column)
			}

			numI := extractNumber(valI)
			numJ := extractNumber(valJ)

			if config.reverse {
				return numI > numJ
			}
			return numI < numJ
		})
	} else if config.human {
		// Human-сортировка (с суффиксами K, M, G)
		sort.Slice(result, func(i, j int) bool {
			valI := result[i]
			valJ := result[j]

			if config.ignoreBlanks {
				valI = trimBlanks(valI)
				valJ = trimBlanks(valJ)
			}

			if config.column > 0 {
				valI = getColumn(valI, config.column)
				valJ = getColumn(valJ, config.column)
			}

			sizeI := parseHumanSize(valI)
			sizeJ := parseHumanSize(valJ)

			if config.reverse {
				return sizeI > sizeJ
			}
			return sizeI < sizeJ
		})
	} else {
		// Обычная текстовая сортировка
		sort.Slice(result, func(i, j int) bool {
			valI := result[i]
			valJ := result[j]

			if config.ignoreBlanks {
				valI = trimBlanks(valI)
				valJ = trimBlanks(valJ)
			}

			if config.column > 0 {
				valI = getColumn(valI, config.column)
				valJ = getColumn(valJ, config.column)
			}

			if config.reverse {
				return valI > valJ
			}
			return valI < valJ
		})
	}

	// Удаление дубликатов
	if config.unique {
		result = uniqueLines(result, config)
	}

	return result
}

// uniqueLines удаляет повторяющиеся строки с учётом конфигурации
func uniqueLines(lines []string, config *ConfigSort) []string {
	if len(lines) <= 1 {
		return lines
	}

	result := make([]string, 0, len(lines))

	// Функция для подготовки строки к сравнению
	prepare := func(s string) string {
		if config.ignoreBlanks {
			s = trimBlanks(s)
		}
		if config.column > 0 {
			s = getColumn(s, config.column)
		}
		if config.month {
			s = extractMonthName(s)
		}
		if config.human {
			return fmt.Sprintf("%f", parseHumanSize(s))
		}
		return s
	}

	result = append(result, lines[0])
	prevPrepared := prepare(lines[0])

	for i := 1; i < len(lines); i++ {
		currentPrepared := prepare(lines[i])
		if currentPrepared != prevPrepared {
			result = append(result, lines[i])
			prevPrepared = currentPrepared
		}
	}

	return result
}

// extractNumber извлекает число из строки
func extractNumber(s string) float64 {
	s = strings.TrimSpace(s)

	num, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return num
	}

	var numStr string
	for _, ch := range s {
		if (ch >= '0' && ch <= '9') || ch == '.' || ch == '-' {
			numStr += string(ch)
		} else if numStr != "" {
			break
		}
	}

	if numStr == "" {
		return 0
	}

	num, _ = strconv.ParseFloat(numStr, 64)
	return num
}

// parseHumanSize парсит размер с суффиксом (1K, 2M, 3G и т.д.)
func parseHumanSize(s string) float64 {
	s = strings.TrimSpace(s)

	var numStr string
	var suffix string

	for i, ch := range s {
		if (ch >= '0' && ch <= '9') || ch == '.' {
			numStr += string(ch)
		} else if ch == ' ' {
			continue
		} else {
			suffix = strings.ToUpper(string(ch))
			// Проверяем остальную часть строки
			remaining := strings.ToUpper(s[i+1:])
			for _, r := range remaining {
				if r >= 'A' && r <= 'Z' {
					suffix += string(r)
				}
			}
			break
		}
	}

	if numStr == "" {
		return extractNumber(s)
	}

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return extractNumber(s)
	}

	switch suffix {
	case "K", "KB":
		return num * 1024
	case "M", "MB":
		return num * 1024 * 1024
	case "G", "GB":
		return num * 1024 * 1024 * 1024
	case "T", "TB":
		return num * 1024 * 1024 * 1024 * 1024
	case "P", "PB":
		return num * 1024 * 1024 * 1024 * 1024 * 1024
	default:
		return num
	}
}

// extractMonthName извлекает название месяца из строки
func extractMonthName(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 3 {
		return s[:3]
	}
	return s
}

// getColumn извлекает указанную колонку из строки
func getColumn(line string, column int) string {
	columns := strings.Split(line, "\t")

	if column-1 < len(columns) {
		return columns[column-1]
	}

	return ""
}

// trimBlanks удаляет хвостовые пробелы
func trimBlanks(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

// isSorted проверяет, отсортированы ли строки
func isSorted(lines []string, config *ConfigSort) bool {
	if len(lines) <= 1 {
		return true
	}

	for i := 1; i < len(lines); i++ {
		valPrev := lines[i-1]
		valCurr := lines[i]

		if config.ignoreBlanks {
			valPrev = trimBlanks(valPrev)
			valCurr = trimBlanks(valCurr)
		}

		if config.column > 0 {
			valPrev = getColumn(valPrev, config.column)
			valCurr = getColumn(valCurr, config.column)
		}

		var less bool

		if config.numeric {
			numPrev := extractNumber(valPrev)
			numCurr := extractNumber(valCurr)
			less = numPrev < numCurr
		} else if config.month {
			monthPrev := monthMap[strings.ToLower(extractMonthName(valPrev))]
			monthCurr := monthMap[strings.ToLower(extractMonthName(valCurr))]
			less = monthPrev < monthCurr
		} else if config.human {
			sizePrev := parseHumanSize(valPrev)
			sizeCurr := parseHumanSize(valCurr)
			less = sizePrev < sizeCurr
		} else {
			less = valPrev < valCurr
		}

		if !less && valPrev != valCurr {
			return false
		}
	}

	return true
}

// readInputData читает данные из STDIN или файла
func readInputData(config *ConfigSort) ([]string, error) {
	if flag.NArg() == 0 {
		return readLines(os.Stdin, config)
	}

	filename := flag.Arg(0)

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла %s: %w", filename, err)
	}
	defer file.Close()

	return readLines(file, config)
}

// readLines читает строки из io.Reader
func readLines(r io.Reader, config *ConfigSort) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var lines []string

	// Увеличиваем буфер для больших строк
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()

		// Если нужно игнорировать пробелы - обрезаем сразу
		if config.ignoreBlanks {
			line = trimBlanks(line)
		}

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
