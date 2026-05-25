package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	go func() {
		for {
			<-sigChan

			fmt.Print("\nminishell> ")
			os.Stdout.Sync()
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("minishell> ")
		os.Stdout.Sync()

		// либо ошибка, либо CTRL + D
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Ошибка чтения: %v\n", err)
			}
			fmt.Println("\nВыход из minishell...")
			break
		}

		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.Contains(line, "|") {
			commands := parsePipeline(line)
			executePipeline(commands)
		} else {
			args := parsing(line)
			if args != nil {
				executeCommand(args)
			}
		}

	}
}

func parsing(line string) []string {
	args := strings.Fields(line)
	if len(args) == 0 {
		return nil
	}

	return args
}

func executeCommand(args []string) {
	commandName := args[0]

	switch commandName {
	case "cd":
		var targetDir string
		var err error

		if len(args) < 2 {
			targetDir, err = os.UserHomeDir()

			if err != nil {
				fmt.Fprintf(os.Stderr, "cd: невозможно получить домашнюю директорию: %v\n", err)
				return
			}
		} else {
			targetDir = args[1]
		}

		err = os.Chdir(targetDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Такой директории не существует: %v\n", err)
		}

	case "pwd":
		dir, err := os.Getwd()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		} else {
			fmt.Println(dir)
		}

	case "echo":
		if len(args) < 2 {
			fmt.Println()
		} else {
			fmt.Println(strings.Join(args[1:], " "))
		}
	case "kill":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Мало аргументов")
			return
		}

		killProcess(args[1])

	case "ps":
		builtinPS()

	default:
		executeExternalCommand(args)
	}

}

func builtinPS() {
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("ps")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка выполнения ps: %v\n", err)
		}
		return
	}

	// 1. Открываем директорию /proc
	files, err := os.ReadDir("/proc")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения /proc: %v\n", err)
		return
	}

	// Печатаем красивую шапку таблицы
	fmt.Printf("%-8s %s\n", "PID", "CMD")
	fmt.Println(strings.Repeat("-", 20))

	// 2. Итерируемся по всему содержимому /proc
	for _, file := range files {
		// Проверяем, что это директория
		if !file.IsDir() {
			continue
		}

		// Проверяем, является ли имя папки числом (PID)
		pidStr := file.Name()
		if _, err := strconv.Atoi(pidStr); err != nil {
			// Если это не число (например, папка /proc/sys), пропускаем её
			continue
		}

		// 3. Формируем путь к файлу 'comm', где лежит имя процесса
		// Например: /proc/1234/comm
		commPath := filepath.Join("/proc", pidStr, "comm")

		// Читаем содержимое файла
		commBytes, err := os.ReadFile(commPath)
		if err != nil {
			// Процесс мог закрыться прямо во время чтения, это нормально — пропускаем
			continue
		}

		// Очищаем имя от лишних переносов строк и пробелов
		cmdName := strings.TrimSpace(string(commBytes))

		// 4. Выводим результат в таблицу
		// %-8s означает выравнивание по левому краю с шириной в 8 символов
		fmt.Printf("%-8s %s\n", pidStr, cmdName)
	}
}

func killProcess(PID string) {
	pid, err := strconv.Atoi(PID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		return
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		return
	}

	err = proc.Signal(syscall.SIGKILL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		return
	}

	fmt.Printf("Процесс %d успешно завершен\n", pid)
}

func executeExternalCommand(args []string) {
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		return
	}
}

func parsePipeline(line string) [][]string {
	var commands [][]string

	stages := strings.Split(line, "|")

	for _, stage := range stages {
		stage = strings.TrimSpace(stage)
		if stage == "" {
			continue
		}
		args := strings.Fields(stage)
		if len(args) > 0 {
			commands = append(commands, args)
		}
	}
	return commands
}

func executePipeline(commands [][]string) {
	numCmds := len(commands)
	if numCmds == 0 {
		return
	}

	// Создаем слайс для команд, чтобы управлять ими в цикле
	cmds := make([]*exec.Cmd, numCmds)

	// Переменная для хранения "выхода" предыдущей команды, который станет "входом" для следующей
	var nextStdin io.ReadCloser

	// 1. Конфигурируем все команды и связываем их шлангами (пайпами)
	for i, args := range commands {
		cmds[i] = exec.Command(args[0], args[1:]...)

		// Если это НЕ первая команда, значит, у нас есть выход от предыдущей.
		// Втыкаем его в Stdin текущей команды.
		if i > 0 {
			cmds[i].Stdin = nextStdin
		} else {
			// Первая команда читает обычный ввод шелла
			cmds[i].Stdin = os.Stdin
		}

		// Если это НЕ последняя команда, ей нужен пайп для передачи данных следующей
		if i < numCmds-1 {
			pr, pw := io.Pipe()
			cmds[i].Stdout = pw // Текущая команда пишет в этот пайп
			cmds[i].Stderr = os.Stderr
			nextStdin = pr // Следующая команда будет читать из этого пайпа
		} else {
			// Последняя команда выводит результат на экран
			cmds[i].Stdout = os.Stdout
			cmds[i].Stderr = os.Stderr
		}
	}

	// 2. Запускаем ВСЕ процессы асинхронно через .Start()
	// (Если запустить через .Run(), шелл зависнет на первой команде, ожидая её конца)
	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка запуска команды: %v\n", err)
			return
		}
	}

	// 3. Магия закрытия ресурсов и синхронизации
	// Нам нужно запустить горутину, которая будет закрывать пайпы по цепочке,
	// когда предыдущие команды завершают работу. Иначе следующие команды будут вечно ждать EOF.
	for i := 0; i < numCmds-1; i++ {
		go func(index int) {
			cmds[index].Wait() // Ждем завершения конкретной команды
			// Закрываем писатель пайпа этой команды.
			// Это пошлет сигнал EOF (конец файла) для Stdin следующей команды.
			if wc, ok := cmds[index].Stdout.(io.Closer); ok {
				wc.Close()
			}
		}(i)
	}

	// 4. Главный поток шелла блокируется и ждет завершения ТОЛЬКО ПОСЛЕДНЕЙ команды
	if err := cmds[numCmds-1].Wait(); err != nil {
		// Ошибки выполнения (например, grep ничего не нашел) можно не выводить как краш шелла
	}
}
