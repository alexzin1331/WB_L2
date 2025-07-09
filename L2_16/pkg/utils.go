package pkg

import (
	"net/url"
	"path/filepath"
	"strings"
)

// ResolveURL преобразует относительный путь href в абсолютный URL,
// используя базовый URL base.
func ResolveURL(href string, base *url.URL) string {
	ref, err := url.Parse(href)
	if err != nil {
		return ""
	}
	return base.ResolveReference(ref).String()
}

// IsAsset проверяет, является ли ссылка ссылкой на внешний ресурс:
// скрипт, стиль, изображение, шрифт и т.д.
func IsAsset(u string) bool {
	return strings.HasSuffix(u, ".js") || strings.HasSuffix(u, ".css") ||
		strings.HasSuffix(u, ".png") || strings.HasSuffix(u, ".jpg") ||
		strings.HasSuffix(u, ".jpeg") || strings.HasSuffix(u, ".gif") ||
		strings.HasSuffix(u, ".svg") || strings.HasSuffix(u, ".woff") ||
		strings.HasSuffix(u, ".ttf") || strings.HasSuffix(u, ".ico")
}

// URLToFilePath преобразует URL в путь к локальному файлу,
// с учётом базового URL и корневой директории outDir.
// Например, https://example.com/path/ → mirror/example.com/path/index.html
func URLToFilePath(u string, base *url.URL, outDir string) string {
	parsed, _ := url.Parse(u)
	path := parsed.Path
	if path == "" || strings.HasSuffix(path, "/") {
		path += "index.html"
	}
	return filepath.Join(outDir, parsed.Host, path)
}

// SameDomain проверяет, принадлежит ли URL тому же домену,
// что и базовый URL (без учёта схемы).
func SameDomain(u string, base *url.URL) bool {
	parsed, err := url.Parse(u)
	if err != nil {
		return false
	}
	return parsed.Host == base.Host
}
