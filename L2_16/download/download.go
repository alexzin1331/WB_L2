package downloader

import (
	"WB_L2/L2_16/parser"
	"WB_L2/L2_16/pkg"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// для посещенных адресов, мьютекс для конкурентной записи
type Visits struct {
	m map[string]bool
	sync.RWMutex
}

func initVisits() *Visits {
	return &Visits{
		m: make(map[string]bool),
	}
}

func Start(startURL string, maxDepth int, out string, workers int) error {
	//инициализируем структуру
	visited := initVisits()

	//готовим ссылку к использованию
	base, err := url.Parse(startURL)
	if err != nil {
		return err
	}

	//создаем заданную директорию
	os.MkdirAll(out, os.ModePerm)
	sem := make(chan struct{}, workers)
	wg := &sync.WaitGroup{}

	var recSearch func(string, int)
	recSearch = func(u string, depth int) {
		//проверка на глубину рекурсии
		if depth <= 0 {
			return
		}

		//блокируем чтение
		visited.RLock()
		if visited.m[u] {
			visited.RUnlock()
			return
		}
		visited.RUnlock()

		//блокируем запись
		visited.Lock()
		visited.m[u] = true
		visited.Unlock()

		//создаем семафор
		sem <- struct{}{}
		wg.Add(1)

		//создаем горутину
		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			//возвращаем исходный HTML, список ресурсов и список ссылок
			html, resources, links, err := parser.FetchAndParse(u)
			if err != nil {
				fmt.Println("Ошибка загрузки:", u, err)
				return
			}

			//преобразует URL в путь к локальному файлу
			localPath := pkg.URLToFilePath(u, base, out)
			//создаем локально директорию и записываем туда html
			os.MkdirAll(filepath.Dir(localPath), os.ModePerm)
			os.WriteFile(localPath, []byte(html), 0644)

			//скачиваем каждый ресурс
			for _, res := range resources {
				go downloadStuff(visited, res, base, out)
			}

			//проходимся по ссылкам на данном уровне рекурсии и если ссылка
			//имеет ту же корневую директорию, что и изначальная, то запускаем рекурсию внутри нее
			for _, link := range links {
				if pkg.SameDomain(link, base) {
					recSearch(link, depth-1)
				}
			}
		}()
	}

	//запускаемся относительно начальной ссылки
	recSearch(startURL, maxDepth)
	wg.Wait()
	return nil
}

// скачиваем ресурс
func downloadStuff(visited *Visits, resURL string, baseUrl *url.URL, out string) {
	//блокируем чтение
	visited.RLock()
	if visited.m[resURL] {
		visited.RUnlock()
		return
	}
	visited.RUnlock()

	//блокируем запись
	visited.Lock()
	visited.m[resURL] = true
	visited.Unlock()

	client := &http.Client{Timeout: 10 * time.Second}
	//формируем запрос (контекст для таймаута)
	req, _ := http.NewRequestWithContext(context.Background(), "GET", resURL, nil)
	//делаем запрос
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return
	}
	defer resp.Body.Close()

	//создаем директорию и записываем туда файл с телом ответа на запрос
	localPath := pkg.URLToFilePath(resURL, baseUrl, out)
	os.MkdirAll(filepath.Dir(localPath), os.ModePerm)
	body, err := io.ReadAll(resp.Body)
	if err == nil {
		os.WriteFile(localPath, body, 0644)
	}
}
