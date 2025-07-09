package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

func main() {
	// канал для обработки сигнала SIGINT (Ctrl+C) echo "one\ntwo\nthree" | wc -l > 1.txt
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	reader := bufio.NewReader(os.Stdin)

	for {
		// выводим приглашение командной строки
		fmt.Print("$ ")

		// читаем ввод пользователя до символа новой строки
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// если получен EOF (Ctrl+D) - завершаем работу
				fmt.Println("\nExit")
				os.Exit(0)
			}
			// в случае других ошибок ввода выводим сообщение
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			continue
		}

		// удаляем лишние пробелы и символы новой строки
		input = strings.TrimSpace(input)
		// пропускаем пустые строки
		if input == "" {
			continue
		}

		// запускаем обработку команды в отдельной горутине
		go func() {
			// ожидаем либо сигнал прерывания, либо продолжаем выполнение
			select {
			case <-signalChan:
				// при получении SIGINT выводим сообщение
				fmt.Println("\nInterrupted")
				return
			default:
				// если сигнала нет, продолжаем выполнение
			}

			// разбиваем команду на части по символу конвейера | и выполняем конвейер команд
			commands := strings.Split(input, "|")
			executePipeline(commands)
		}()
	}
}

// выполняет конвейер команд
func executePipeline(commands []string) {
	var prevCmd *exec.Cmd        // предыдущая команда в конвейере
	var prevOutput io.ReadCloser // выходной поток предыдущей команды

	// обрабатываем каждую команду в конвейере
	for i, cmd := range commands {
		// удаляем лишние пробелы и разбиваем команду на аргументы
		cmd = strings.TrimSpace(cmd)
		args := parseCommand(cmd) // новая функция для обработки аргументов
		if len(args) == 0 {
			continue
		}

		// проверяем, является ли команда встроенной
		if handleBuiltinCommand(args) {
			continue
		}

		// создаем команду для выполнения
		cmdPath, err := CommandPath(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		cmdExec := exec.Command(cmdPath, args[1:]...)

		// обработка редиректов
		cmdExec = handleRedirections(cmdExec, args)

		//ввод
		// для первой команды - stdin консоли
		// для следующих - вывод предыдущей команды
		if i > 0 {
			cmdExec.Stdin = prevOutput
		} else {
			cmdExec.Stdin = os.Stdin
		}
		//вывод
		// для последней команды - stdout консоли
		// для остальных - создаем pipe для связи со следующей командой
		if i < len(commands)-1 {
			prevOutput, _ = cmdExec.StdoutPipe()
		} else {
			cmdExec.Stdout = os.Stdout
		}

		// ошибки выводим в stderr консоли
		cmdExec.Stderr = os.Stderr

		// запускаем команду
		err = cmdExec.Start()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error starting command:", err)
			return
		}

		// ожидаем завершения предыдущей команды (если она есть)
		if prevCmd != nil {
			prevCmd.Wait()
		}
		prevCmd = cmdExec
	}

	// Ожидаем завершения последней команды в конвейере
	if prevCmd != nil {
		prevCmd.Wait()
	}
}

// parseCommand обрабатывает команду, заменяет переменные окружения и возвращает аргументы
func parseCommand(cmd string) []string {
	args := strings.Fields(cmd)
	for i, arg := range args {
		// Заменяем $VAR на значение переменной окружения
		if strings.HasPrefix(arg, "$") {
			varName := strings.TrimPrefix(arg, "$")
			args[i] = os.Getenv(varName)
		}
	}
	return args
}

// функция для работы внешних команд
func CommandPath(name string) (string, error) {
	// проверяем абсолютный путь или путь относительно текущей директории
	if strings.Contains(name, "/") {
		if _, err := os.Stat(name); err == nil {
			return name, nil
		}
		return "", fmt.Errorf("error resolve Command Path")
	}

	// ищем в PATH
	for _, dir := range strings.Split(os.Getenv("PATH"), ":") {
		fullPath := filepath.Join(dir, name)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath, nil
		}
	}

	return "", fmt.Errorf("command not found: %s", name)
}

// handleRedirections обрабатывает редиректы ввода/вывода
func handleRedirections(cmd *exec.Cmd, args []string) *exec.Cmd {
	var newArgs []string
	var inputFile, outputFile string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "<":
			if i+1 < len(args) {
				inputFile = args[i+1]
				i++
			}
		case ">":
			if i+1 < len(args) {
				outputFile = args[i+1]
				i++
			}
		default:
			newArgs = append(newArgs, arg)
		}
	}

	// Устанавливаем редиректы
	if inputFile != "" {
		file, err := os.Open(inputFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error opening input file:", err)
			return cmd
		}
		cmd.Stdin = file
	}

	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating output file:", err)
			return cmd
		}
		cmd.Stdout = file
	}

	cmd.Args = newArgs
	if len(newArgs) > 0 {
		cmd.Path = newArgs[0]
	}

	return cmd
}

// обрабатывает встроенные команды
// возвращает true, если команда была встроенной и обработана
func handleBuiltinCommand(args []string) bool {
	switch args[0] {
	case "cd":
		// команда смены директории
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "cd: missing argument")
			return true
		}
		err := os.Chdir(args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "cd:", err)
		}
		return true

	case "pwd":
		// команда вывода текущей директории
		dir, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "pwd:", err)
			return true
		}
		fmt.Println(dir)
		return true

	case "echo":
		// команда вывода аргументов
		fmt.Println(strings.Join(args[1:], " "))
		return true

	case "kill":
		// команда завершения процесса
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "kill: missing PID")
			return true
		}
		pid := args[1]
		// используем системную команду kill
		cmd := exec.Command("kill", pid)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, "kill:", err)
		}
		return true

	case "ps":
		// команда вывода списка процессов
		// используем системную команду ps
		cmd := exec.Command("ps", "aux")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, "ps:", err)
		}
		return true
	}

	// если команда не встроенная, возвращаем false
	return false
}
