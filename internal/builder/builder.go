package builder

import (
	"fmt"
	"sort"
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
	Format    Format
	Header    string
	Footer    string
	ShowTree  bool
	TreeTitle string
}

func Build(files []scanner.File, opts Options) string {
	var b strings.Builder

	if opts.Header != "" {
		b.WriteString(opts.Header)
		b.WriteString("\n\n")
	}

	if opts.ShowTree && len(files) > 0 {
		title := opts.TreeTitle
		if title == "" {
			title = "Project structure"
		}
		writeTree(&b, files, title, opts.Format)
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

func writeTree(b *strings.Builder, files []scanner.File, title string, format Format) {
	paths := make([]string, 0, len(files))
	for _, f := range files {
		paths = append(paths, f.Path)
	}
	sort.Strings(paths)

	tree := renderTree(paths)

	switch format {
	case FormatMarkdown:
		fmt.Fprintf(b, "## %s\n\n```\n%s```\n", title, tree)
	case FormatXML:
		fmt.Fprintf(b, "<structure>\n%s</structure>\n", tree)
	default:
		fmt.Fprintf(b, "===== %s =====\n\n%s", title, tree)
	}
}

func renderTree(paths []string) string {
	root := &node{children: map[string]*node{}}

	for _, p := range paths {
		parts := strings.Split(p, "/")
		cur := root
		for i, part := range parts {
			if cur.children == nil {
				cur.children = map[string]*node{}
			}
			child, ok := cur.children[part]
			if !ok {
				child = &node{name: part, children: map[string]*node{}}
				cur.children[part] = child
			}
			if i == len(parts)-1 {
				child.isFile = true
			}
			cur = child
		}
	}

	var b strings.Builder
	renderNode(&b, root, "")
	return b.String()
}

type node struct {
	name     string
	isFile   bool
	children map[string]*node
}

func renderNode(b *strings.Builder, n *node, prefix string) {
	names := make([]string, 0, len(n.children))
	for name := range n.children {
		names = append(names, name)
	}
	sort.Strings(names)

	for i, name := range names {
		child := n.children[name]
		last := i == len(names)-1

		connector := "├── "
		nextPrefix := prefix + "│   "
		if last {
			connector = "└── "
			nextPrefix = prefix + "    "
		}

		b.WriteString(prefix)
		b.WriteString(connector)
		b.WriteString(name)
		if !child.isFile {
			b.WriteString("/")
		}
		b.WriteString("\n")

		if len(child.children) > 0 {
			renderNode(b, child, nextPrefix)
		}
	}
}

func writePlain(b *strings.Builder, f scanner.File) {
	fmt.Fprintf(b, "\n===== %s =====\n\n", f.Path)
	b.Write(f.Content)
	b.WriteString("\n")
}

func writeMarkdown(b *strings.Builder, f scanner.File) {
	lang := langFromExt(f.Path)
	if lang == "" {
		lang = "text"
	}
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
	switch strings.ToLower(path[idx+1:]) {
	case "go":
		return "go"
	case "py":
		return "python"
	case "js", "mjs", "cjs":
		return "javascript"
	case "jsx":
		return "jsx"
	case "ts":
		return "typescript"
	case "tsx":
		return "tsx"
	case "rs":
		return "rust"
	case "java":
		return "java"
	case "kt", "kts":
		return "kotlin"
	case "swift":
		return "swift"
	case "c":
		return "c"
	case "h":
		return "c"
	case "cpp", "cc", "cxx":
		return "cpp"
	case "hpp", "hh", "hxx":
		return "cpp"
	case "cs":
		return "csharp"
	case "rb":
		return "ruby"
	case "php":
		return "php"
	case "pl":
		return "perl"
	case "lua":
		return "lua"
	case "sh", "bash", "zsh":
		return "bash"
	case "ps1":
		return "powershell"
	case "bat", "cmd":
		return "batch"
	case "sql":
		return "sql"
	case "html", "htm":
		return "html"
	case "css":
		return "css"
	case "scss":
		return "scss"
	case "sass":
		return "sass"
	case "less":
		return "less"
	case "vue":
		return "vue"
	case "svelte":
		return "svelte"
	case "astro":
		return "astro"
	case "md", "markdown":
		return "markdown"
	case "rst":
		return "rst"
	case "txt":
		return "text"
	case "json":
		return "json"
	case "json5":
		return "json5"
	case "yaml", "yml":
		return "yaml"
	case "toml":
		return "toml"
	case "ini":
		return "ini"
	case "xml":
		return "xml"
	case "svg":
		return "xml"
	case "csv":
		return "csv"
	case "tsv":
		return "tsv"
	case "env":
		return "dotenv"
	case "dockerfile":
		return "dockerfile"
	case "makefile":
		return "makefile"
	case "gradle":
		return "gradle"
	case "proto":
		return "protobuf"
	case "graphql", "gql":
		return "graphql"
	case "tf":
		return "hcl"
	case "tfvars":
		return "hcl"
	case "docx":
		return "text"
	default:
		return path[idx+1:]
	}
}
