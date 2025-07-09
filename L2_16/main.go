package main

import (
	"WB_L2/L2_16/download"
	"flag"
	"fmt"
	"log"
)

func main() {
	// установим флаги и парсим их
	url := flag.String("url", "", "web URL")
	depth := flag.Int("depth", 2, "Recursion depth")
	out := flag.String("out", "mirror", "Directory to save")
	workers := flag.Int("workers", 4, "Number of parallel downloads")
	flag.Parse()

	if *url == "" {
		log.Fatal("Укажите URL с помощью флага -url")
	}

	// начинаем загрузку
	err := downloader.Start(*url, *depth, *out, *workers)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	fmt.Println("Загрузка завершена.")
}
