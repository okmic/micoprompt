package scanner

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Options struct {
	Dir          string
	Extensions   []string
	IgnoreDirs   []string
	IgnoreFiles  []string
	MaxFileSize  int64
	MaxTotal     int64
	UseGitignore bool
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

		for _, ig := range opts.IgnoreFiles {
			if name == ig {
				res.Skipped++
				return nil
			}
		}

		if len(exts) > 0 {
			ext := strings.ToLower(filepath.Ext(path))
			if !contains(exts, ext) {
				res.Skipped++
				return nil
			}
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

		if isBinary(data) {
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
