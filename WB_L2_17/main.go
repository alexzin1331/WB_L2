package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// парсим аргументы командной строки
	timeout := flag.Int("timeout", 10, "connection timeout (s)")
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("usage: telnet <host> <port> [--timeout <seconds>]")
		os.Exit(1)
	}

	host := args[0]
	port := args[1]
	address := net.JoinHostPort(host, port)

	// устанавливаем соединение с таймаутом
	conn, err := net.DialTimeout("tcp", address, time.Duration(*timeout)*time.Second)
	if err != nil {
		log.Printf("error connecting to %s: %v\n", address, err)
		os.Exit(1)
	}
	defer conn.Close()

	log.Printf("connected to %s\n", address)

	// каналы для управления завершением
	done := make(chan struct{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// читаем из сокета и выводим в STDOUT
	go func() {
		reader := bufio.NewReader(conn)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					fmt.Println("\nconnection closed by remote host")
				} else {
					fmt.Printf("\nread error: %v\n", err)
				}
				close(done)
				return
			}
			fmt.Print(line)
		}
	}()

	// читаем из STDIN и пишем в сокет
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text() + "\n"
			_, err := conn.Write([]byte(text))
			if err != nil {
				fmt.Printf("write error: %v\n", err)
				close(done)
				return
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("stdin read error: %v\n", err)
		}
		close(done)
	}()

	// ждем сигнала завершения или закрытия соединения
	select {
	case <-done:
	case <-sigChan:
		fmt.Println("\nreceived interrupt signal, closing connection...")
	}
}
