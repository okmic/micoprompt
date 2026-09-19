package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Mico/micoprompt/internal/builder"
	"github.com/Mico/micoprompt/internal/scanner"
	"github.com/Mico/micoprompt/pkg/tokenizer"
)

const version = "0.1.0"

func main() {
	var (
		dir         = flag.String("dir", ".", "папка для обхода")
		out         = flag.String("out", "prompt.txt", "выходной файл (- для stdout)")
		ext         = flag.String("ext", ".go,.md,.txt,.mod,.yaml,.yml,.json", "расширения через запятую")
		ignore      = flag.String("ignore", ".git,node_modules,vendor,dist,build,.idea,.vscode", "игнорируемые папки")
		ignoreFiles = flag.String("ignore-files", ".env,.DS_Store", "игнорируемые имена файлов")
		maxSize     = flag.Int64("max-size", 100*1024, "макс. размер файла в байтах")
		maxTotal    = flag.Int64("max-total", 0, "макс. общий размер (0 = без лимита)")
		format      = flag.String("format", "markdown", "формат: plain | markdown | xml")
		header      = flag.String("header", "", "заголовок в начале промпта")
		footer      = flag.String("footer", "", "футер в конце промпта")
		showVersion = flag.Bool("version", false, "показать версию")
	)

	flag.Usage = usage
	flag.Parse()

	if *showVersion {
		fmt.Println("micoprompt", version)
		return
	}

	start := time.Now()

	// 1. Сканируем
	res, err := scanner.Scan(scanner.Options{
		Dir:         *dir,
		Extensions:  splitTrim(*ext),
		IgnoreDirs:  splitTrim(*ignore),
		IgnoreFiles: splitTrim(*ignoreFiles),
		MaxFileSize: *maxSize,
		MaxTotal:    *maxTotal,
	})
	if err != nil {
		fatal("ошибка сканирования: %v", err)
	}

	prompt := builder.Build(res.Files, builder.Options{
		Format: builder.Format(*format),
		Header: *header,
		Footer: *footer,
	})

	if *out == "-" {
		fmt.Print(prompt)
	} else {
		if err := os.WriteFile(*out, []byte(prompt), 0644); err != nil {
			fatal("ошибка записи: %v", err)
		}
	}

	tokens := tokenizer.Estimate(prompt)
	fmt.Printf("✅ Готово за %s\n", time.Since(start).Round(time.Millisecond))
	fmt.Printf("   Файлов:    %d\n", len(res.Files))
	fmt.Printf("   Пропущено: %d\n", res.Skipped)
	fmt.Printf("   Символов:  %d\n", len(prompt))
	fmt.Printf("   Токенов:   ~%d (оценка)\n", tokens)
	fmt.Printf("   Размер:    %.1f KB\n", float64(len(prompt))/1024)
	if *out != "-" {
		fmt.Printf("   Файл:      %s\n", *out)
	}
}

func splitTrim(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := parts[:0]
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "❌ "+format+"\n", args...)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, `micoprompt — собирает файлы проекта в один промпт для ИИ

Использование:
  micoprompt [флаги]

Примеры:
  micoprompt -dir ./src -out prompt.txt
  micoprompt -dir ./src -format xml -out - | pbcopy
  micoprompt -dir ./src -ext .go,.md -ignore .git,vendor

Флаги:
`)
	flag.PrintDefaults()
}
