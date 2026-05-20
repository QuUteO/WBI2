package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

type Config struct {
	after      int
	before     int
	context    int
	count      bool
	ignoreCase bool
	invert     bool
	fixed      bool
	lineNum    bool
	pattern    string
}

func main() {
	config := parseFlags()

	var lines []string

	if flag.NArg() > 1 {
		fileName := flag.Arg(1)
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
	}

	result := grep(lines, config)
	for _, line := range result {
		fmt.Println(line)
	}
}

func parseFlags() Config {
	var config Config

	flag.IntVar(&config.after, "A", 0, "")
	flag.IntVar(&config.before, "B", 0, "")
	flag.IntVar(&config.context, "C", 0, "")
	flag.BoolVar(&config.count, "c", false, "")
	flag.BoolVar(&config.ignoreCase, "i", false, "")
	flag.BoolVar(&config.invert, "v", false, "")
	flag.BoolVar(&config.fixed, "F", false, "")
	flag.BoolVar(&config.lineNum, "n", false, "")

	flag.Parse()

	if flag.NArg() > 0 {
		config.pattern = flag.Arg(0)
	}

	if config.context > 0 {
		config.before = config.context
		config.after = config.context
	}

	return config
}

func grep(lines []string, config Config) []string {
	if len(lines) == 0 {
		if config.count {
			return []string{"0"}
		}
		return []string{}
	}

	// Раскрываем context внутри grep, так как тесты вызывают grep напрямую в обход parseFlags
	if config.context > 0 {
		config.before = config.context
		config.after = config.context
	}

	// Находим все совпадающие строки
	matchLines := make(map[int]bool)

	// Подготовка regexp если нужно
	var re *regexp.Regexp
	if !config.fixed {
		pattern := config.pattern
		if config.ignoreCase {
			pattern = "(?i)" + pattern
		}
		var err error
		re, err = regexp.Compile(pattern)
		if err != nil {
			return []string{}
		}
	}

	for i, line := range lines {
		var matched bool

		if config.fixed {
			// Для -F (fixed) ищем строгое совпадение всей строки, а не подстроки
			if config.ignoreCase {
				matched = strings.ToLower(line) == strings.ToLower(config.pattern)
			} else {
				matched = line == config.pattern
			}
		} else {
			matched = re.MatchString(line)
		}

		if config.invert {
			matched = !matched
		}

		if matched {
			matchLines[i] = true
		}
	}

	// Если нужен только подсчет
	if config.count {
		return []string{fmt.Sprintf("%d", len(matchLines))}
	}

	// Если нет совпадений
	if len(matchLines) == 0 {
		return []string{}
	}

	// Если не нужно выводить контекст
	if config.before == 0 && config.after == 0 {
		result := make([]string, 0)
		for i := 0; i < len(lines); i++ {
			if matchLines[i] {
				if config.lineNum {
					result = append(result, fmt.Sprintf("%d:%s", i+1, lines[i]))
				} else {
					result = append(result, lines[i])
				}
			}
		}
		return result
	}

	// Вывод с контекстом
	type interval struct {
		start int
		end   int
	}

	intervals := make([]interval, 0)

	// Сортируем совпадающие строки
	matchSorted := make([]int, 0, len(matchLines))
	for k := range matchLines {
		matchSorted = append(matchSorted, k)
	}
	sort.Ints(matchSorted)

	// Объединяем пересекающиеся интервалы
	for _, idx := range matchSorted {
		start := idx - config.before
		if start < 0 {
			start = 0
		}
		end := idx + config.after
		if end >= len(lines) {
			end = len(lines) - 1
		}

		// Объединяем, только если интервалы НАКЛАДЫВАЮТСЯ друг на друга (start <= end).
		// Если они идут встык (start == end + 1), они не объединяются, чтобы остался разделитель '--'.
		if len(intervals) > 0 && start <= intervals[len(intervals)-1].end {
			if end > intervals[len(intervals)-1].end {
				intervals[len(intervals)-1].end = end
			}
		} else {
			intervals = append(intervals, interval{start, end})
		}
	}

	// Формируем результат
	result := make([]string, 0)
	for i, inter := range intervals {
		// Добавляем разделитель между изолированными интервалами
		if i > 0 {
			result = append(result, "--")
		}

		// Выводим строки интервала
		for j := inter.start; j <= inter.end; j++ {
			var output string
			if config.lineNum {
				output = fmt.Sprintf("%d:%s", j+1, lines[j])
			} else {
				output = lines[j]
			}
			result = append(result, output)
		}
	}

	return result
}
