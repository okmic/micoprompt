package builder

import (
	"fmt"
	"strings"

	"github.com/Mico/micoprompt/internal/scanner"
)

type Format string

const (
	FormatPlain    Format = "plain"
	FormatMarkdown Format = "markdown"
	FormatXML      Format = "xml"
)

type Options struct {
	Format Format
	Header string
	Footer string
}

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
	case "jsx":
		return "react-javascript"
	case "ts":
		return "typescript"
	case "tsx":
		return "react-typescript"
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
