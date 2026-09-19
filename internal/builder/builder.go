package builder

import (
	"fmt"
	"strings"

	"github.com/Mico/micoprompt/internal/scanner"
)

// Format — формат вывода.
type Format string

const (
	FormatPlain    Format = "plain"    // ===== FILE =====
	FormatMarkdown Format = "markdown" // ## FILE + ```code```
	FormatXML      Format = "xml"      // <file path="...">...</file>
)

// Options — настройки сборки.
type Options struct {
	Format Format
	Header string // произвольный заголовок в начале промпта
	Footer string // произвольный футер
}

// Build собирает все файлы в один промпт.
func Build(files []scanner.File, opts Options) string {
	var b strings.Builder

	if opts.Header != "" {
		b.WriteString(opts.Header)
		b.WriteString("\n\n")
	}

	for _, f := range files {
		switch opts.Format {
		case FormatMarkdown:
			writeMarkdown(&b, f)
		case FormatXML:
			writeXML(&b, f)
		default:
			writePlain(&b, f)
		}
	}

	if opts.Footer != "" {
		b.WriteString("\n")
		b.WriteString(opts.Footer)
	}

	return b.String()
}

func writePlain(b *strings.Builder, f scanner.File) {
	fmt.Fprintf(b, "\n===== %s =====\n\n", f.Path)
	b.Write(f.Content)
	b.WriteString("\n")
}

func writeMarkdown(b *strings.Builder, f scanner.File) {
	lang := langFromExt(f.Path)
	fmt.Fprintf(b, "\n## %s\n\n```%s\n", f.Path, lang)
	b.Write(f.Content)
	b.WriteString("\n```\n")
}

func writeXML(b *strings.Builder, f scanner.File) {
	fmt.Fprintf(b, "\n<file path=%q>\n", f.Path)
	b.Write(f.Content)
	b.WriteString("\n</file>\n")
}

func langFromExt(path string) string {
	idx := strings.LastIndex(path, ".")
	if idx == -1 {
		return ""
	}
	switch path[idx+1:] {
	case "go":
		return "go"
	case "py":
		return "python"
	case "js":
		return "javascript"
	case "ts":
		return "typescript"
	case "rs":
		return "rust"
	case "md":
		return "markdown"
	case "json":
		return "json"
	case "yaml", "yml":
		return "yaml"
	default:
		return path[idx+1:]
	}
}
