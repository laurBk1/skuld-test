package fileutil

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexmullins/zip"
)

// AppendFile adaugă un rând la sfârșitul fișierului
func AppendFile(path string, line string) {
	file, _ := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	defer file.Close()
	file.WriteString(line + "\n")
}

// Tree generează un view tip arbore al directorului
func Tree(path string, prefix string, isFirstDir ...bool) string {
	var sb strings.Builder

	files, _ := ioutil.ReadDir(path)
	for i, file := range files {
		isLast := i == len(files)-1
		var pointer string
		if isLast {
			pointer = prefix + "└── "
		} else {
			pointer = prefix + "├── "
		}
		if isFirstDir == nil {
			pointer = prefix
		}
		if file.IsDir() {
			fmt.Fprintf(&sb, "%s📂 - %s\n", pointer, file.Name())
			if isLast {
				sb.WriteString(Tree(filepath.Join(path, file.Name()), prefix+"    ", false))
			} else {
				sb.WriteString(Tree(filepath.Join(path, file.Name()), prefix+"│   ", false))
			}
		} else {
			sizeKB := float64(file.Size()) / 1024
			if sizeKB < 1024 {
				fmt.Fprintf(&sb, "%s📄 - %s (%.2f KB)\n", pointer, file.Name(), sizeKB)
			} else {
				fmt.Fprintf(&sb, "%s📄 - %s (%.2f MB)\n", pointer, file.Name(), sizeKB/1024)
			}
		}
	}

	tree := sb.String()
	if len(tree) > 3000 {
		lines := strings.Split(tree, "\n")
		if len(lines) > 100 {
			tree = strings.Join(lines[:100], "\n") + "\n... (truncated - too many files)"
		}
	}
	return tree
}

// Zip comprimă un director într-un fișier zip
func Zip(dirPath string, zipName string) error {
	zipFile, err := os.Create(zipName)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(dirPath, filePath)
		if err != nil {
			return err
		}

		zipEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(zipEntry, file)
		return err
	})

	return err
}

// GetDirectorySize returnează dimensiunea totală a directorului
func GetDirectorySize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

// CountFiles numără toate fișierele dintr-un director
func CountFiles(path string) (int, error) {
	var count int
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	return count, err
}

// Exists verifică dacă un fișier sau folder există
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsDir verifică dacă path-ul este un director
func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

// ReadFile citește întregul fișier ca string
func ReadFile(path string) (string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ReadLines citește fișierul linie cu linie
func ReadLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := make([]string, 0)
	buf := bufio.NewReader(f)

	for {
		line, _, err := buf.ReadLine()
		l := string(line)
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		result = append(result, l)
	}

	return result, nil
}

// WriteFile scrie un string într-un fișier
func WriteFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// Copy copiază un fișier sau director
func Copy(src, dst string) error {
	si, err := os.Stat(src)
	if err != nil {
		return err
	}

	if si.IsDir() {
		return CopyDirSafe(src, dst)
	} else {
		return CopyFileSafe(src, dst)
	}
}

// CopyFileSafe copiază un fișier fără a șterge fișiere existente
func CopyFileSafe(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	os.MkdirAll(filepath.Dir(dst), os.ModePerm)

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	if err := out.Sync(); err != nil {
		return err
	}

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, si.Mode())
}

// CopyDirSafe copiază un director complet fără a șterge folderul destinație
func CopyDirSafe(src, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	if !Exists(dst) {
		if err := os.MkdirAll(dst, si.Mode()); err != nil {
			return err
		}
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := CopyDirSafe(srcPath, dstPath); err != nil {
				fmt.Println("Skip folder:", srcPath, "Error:", err)
				continue
			}
		} else {
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}
			if err := CopyFileSafe(srcPath, dstPath); err != nil {
				fmt.Println("Skip file:", srcPath, "Error:", err)
				continue
			}
		}
	}

	return nil
}

// --- EXPORT PENTRU COMPATIBILITATE BUILD EXISTENT ---
var CopyFile = CopyFileSafe
var CopyDir = CopyDirSafe
