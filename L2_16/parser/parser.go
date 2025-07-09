package parser

import (
	"WB_L2/L2_16/pkg"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// FetchAndParse загружает HTML-страницу по заданному URL, возвращает:
// исходный HTML-код страницы,
// список ресурсов (CSS, JS, изображения и т.д.),
// список ссылок на другие страницы того же домена.
func FetchAndParse(pageURL string) (string, []string, []string, error) {
	// выполняем GET-запрос к странице
	resp, err := http.Get(pageURL)
	if err != nil {
		return "", nil, nil, err
	}
	defer resp.Body.Close()

	// читаем тело ответа в виде байтов
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, nil, err
	}

	// парсим HTML-документ из тела ответа
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return "", nil, nil, err
	}

	// базовый URL для нормализации относительных ссылок
	baseURL, _ := url.Parse(pageURL)

	var resources []string // ссылки на ресурсы (изображения, стили и т.д.)
	var links []string     // ссылки на другие HTML-страницы внутри домена

	// рекурсивная функция обхода DOM-дерева HTML
	var rec func(*html.Node)
	rec = func(node *html.Node) {
		if node.Type == html.ElementNode {
			for _, attr := range node.Attr {
				switch attr.Key {
				case "src", "href":
					ref := attr.Val
					abs := pkg.ResolveURL(ref, baseURL) // преобразуем в абсолютный URL

					//если ref это ресурс, то добавляем
					if pkg.IsAsset(ref) {
						resources = append(resources, abs) // добавляем в список ресурсов
					} else if strings.HasPrefix(abs, baseURL.Scheme+"://"+baseURL.Host) {
						links = append(links, abs) // добавляем в список внутренних ссылок
					}
				}
			}
		}

		// рекурсивно обходим дочерние узлы
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			rec(c)
		}
	}

	// запускаем парсинг с корня документа go run L2_16/main.go -url https://httpbin.org/ -depth 2 -out ./mirror -workers 5
	rec(doc)

	// возвращаем исходный HTML, список ресурсов и список ссылок
	return string(body), resources, links, nil
}
