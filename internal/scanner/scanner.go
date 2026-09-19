package scanner

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Options — настройки сканирования.
type Options struct {
	Dir          string
	Extensions   []string // например: [".go", ".md"]
	IgnoreDirs   []string // например: [".git", "node_modules"]
	IgnoreFiles  []string // конкретные имена
	MaxFileSize  int64    // в байтах
	MaxTotal     int64    // в байтах, 0 = без лимита
	UseGitignore bool     // читать .gitignore
}

// File — найденный файл.
type File struct {
	Path    string // относительный путь
	Content []byte
	Size    int64
}

// Result — итог сканирования.
type Result struct {
	Files     []File
	Skipped   int
	TotalSize int64
}

// Scan обходит папку и возвращает отфильтрованные файлы.
func Scan(opts Options) (*Result, error) {
	res := &Result{}

	// нормализуем расширения (на случай ".go" vs "go")
	exts := make([]string, len(opts.Extensions))
	for i, e := range opts.Extensions {
		if !strings.HasPrefix(e, ".") {
			e = "." + e
		}
		exts[i] = strings.ToLower(e)
	}

	ignoreSet := make(map[string]bool, len(opts.IgnoreDirs))
	for _, d := range opts.IgnoreDirs {
		ignoreSet[d] = true
	}

	err := filepath.WalkDir(opts.Dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // пропускаем недоступные файлы
		}

		name := d.Name()

		// игнорируем папки целиком
		if d.IsDir() {
			if ignoreSet[name] {
				return filepath.SkipDir
			}
			return nil
		}

		// фильтр по имени файла
		for _, ig := range opts.IgnoreFiles {
			if name == ig {
				res.Skipped++
				return nil
			}
		}

		// фильтр по расширению
		if len(exts) > 0 {
			ext := strings.ToLower(filepath.Ext(path))
			if !contains(exts, ext) {
				res.Skipped++
				return nil
			}
		}

		// размер файла
		info, err := d.Info()
		if err != nil {
			res.Skipped++
			return nil
		}
		if opts.MaxFileSize > 0 && info.Size() > opts.MaxFileSize {
			res.Skipped++
			return nil
		}

		// общий лимит
		if opts.MaxTotal > 0 && res.TotalSize+info.Size() > opts.MaxTotal {
			res.Skipped++
			return nil
		}

		// читаем содержимое
		data, err := os.ReadFile(path)
		if err != nil {
			res.Skipped++
			return nil
		}

		// бинарник?
		if isBinary(data) {
			res.Skipped++
			return nil
		}

		// относительный путь для красоты
		rel, err := filepath.Rel(opts.Dir, path)
		if err != nil {
			rel = path
		}

		res.Files = append(res.Files, File{
			Path:    filepath.ToSlash(rel),
			Content: data,
			Size:    info.Size(),
		})
		res.TotalSize += info.Size()
		return nil
	})

	return res, err
}

func contains(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}

// isBinary определяет, бинарный ли файл (по наличию \x00 в первых 512 байтах).
func isBinary(data []byte) bool {
	n := len(data)
	if n > 512 {
		n = 512
	}
	for i := 0; i < n; i++ {
		if data[i] == 0 {
			return true
		}
	}
	return false
}
