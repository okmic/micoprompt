package scanner

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Options struct {
	Dir         string
	Extensions  []string
	IgnoreDirs  []string
	IgnoreFiles []string
	MaxFileSize int64
	MaxTotal    int64
	AllText     bool
}

type File struct {
	Path    string
	Content []byte
	Size    int64
}

type Result struct {
	Files     []File
	Skipped   int
	TotalSize int64
}

var alwaysIgnoreFiles = map[string]bool{
	"micoprompt.txt": true,
}

func Scan(opts Options) (*Result, error) {
	res := &Result{}

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

	userIgnoreFiles := make(map[string]bool, len(opts.IgnoreFiles))
	for _, f := range opts.IgnoreFiles {
		userIgnoreFiles[f] = true
	}

	err := filepath.WalkDir(opts.Dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		name := d.Name()

		if d.IsDir() {
			if ignoreSet[name] {
				return filepath.SkipDir
			}
			return nil
		}

		if alwaysIgnoreFiles[name] || userIgnoreFiles[name] {
			res.Skipped++
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))

		if isAlwaysSkip(ext) {
			res.Skipped++
			return nil
		}

		if !opts.AllText && len(exts) > 0 && !contains(exts, ext) {
			res.Skipped++
			return nil
		}

		info, err := d.Info()
		if err != nil {
			res.Skipped++
			return nil
		}
		if opts.MaxFileSize > 0 && info.Size() > opts.MaxFileSize {
			res.Skipped++
			return nil
		}
		if opts.MaxTotal > 0 && res.TotalSize+info.Size() > opts.MaxTotal {
			res.Skipped++
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			res.Skipped++
			return nil
		}

		if ext == ".docx" {
			text, ok := extractDocxText(data)
			if !ok {
				res.Skipped++
				return nil
			}
			data = []byte(text)
		} else if isBinary(data) {
			res.Skipped++
			return nil
		}

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

func isAlwaysSkip(ext string) bool {
	switch ext {
	case ".exe", ".dll", ".so", ".dylib", ".bin",
		".png", ".jpg", ".jpeg", ".gif", ".webp", ".bmp", ".ico",
		".mp3", ".mp4", ".wav", ".avi", ".mov", ".mkv", ".webm",
		".zip", ".tar", ".gz", ".bz2", ".xz", ".7z", ".rar",
		".pdf", ".doc", ".xls", ".ppt",
		".woff", ".woff2", ".ttf", ".otf", ".eot",
		".lock", ".sum",
		".pyc", ".class", ".o", ".a", ".obj":
		return true
	}
	return false
}

func isBinary(data []byte) bool {
	n := len(data)
	if n > 8000 {
		n = 8000
	}
	if n == 0 {
		return false
	}

	for i := 0; i < n; i++ {
		if data[i] == 0 {
			return true
		}
	}

	odd := 0
	for i := 0; i < n; i++ {
		b := data[i]
		if b < 0x09 || (b > 0x0d && b < 0x20) {
			odd++
		}
	}
	return float64(odd)/float64(n) > 0.30
}

func extractDocxText(data []byte) (string, bool) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", false
	}

	var docXML io.ReadCloser
	for _, f := range zr.File {
		if f.Name == "word/document.xml" {
			docXML, err = f.Open()
			if err != nil {
				return "", false
			}
			break
		}
	}
	if docXML == nil {
		return "", false
	}
	defer docXML.Close()

	return parseDocumentXML(docXML)
}

func parseDocumentXML(r io.Reader) (string, bool) {
	dec := xml.NewDecoder(r)
	var b strings.Builder

	inText := false

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", false
		}

		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "p":
				b.WriteString("\n")
			case "t":
				inText = true
			case "tab":
				b.WriteString("\t")
			case "br":
				b.WriteString("\n")
			}
		case xml.EndElement:
			if t.Name.Local == "t" {
				inText = false
			}
		case xml.CharData:
			if inText {
				b.Write(t)
			}
		}
	}

	return b.String(), true
}
