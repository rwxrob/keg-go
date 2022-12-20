package keg

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func notExists(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

// stringify first looks for string, []byte, []rune, and io.Reader types
// and if matched returns a string with their content and the string
// type.
//
// String converts whatever remains to that types fmt.Sprintf("%v")
// string version (but avoids calling it if possible). Be sure you use
// things with consistent string representations. Keep in mind that this
// is extremely high level for rapid tooling and prototyping.
func stringify(in any) string {
	switch v := in.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case []rune:
		return string(v)
	case io.Reader:
		buf, err := io.ReadAll(v)
		if err != nil {
			log.Println(err)
		}
		return string(buf)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func fetch(url string) ([]byte, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}

	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		return body, ErrFetch{resp}
	}

	return body, nil
}

// DirIsNode returns true if the filepath.Base name of the passed path
// is a valid positive integer (including 0).
func DirIsNode(path string) bool {
	name := filepath.Base(path)
	if i, err := strconv.Atoi(name); i < 0 || err != nil {
		return false
	}
	return true
}

// NodeDirs returns all the directory entries within the target directory
// that have integer names. The lowest integer and highest integer
// values are also returned. Only positive integers are checked. This is
// useful when using directory names as database-friendly unique primary
// keys for other file system content. An empty slice with -1 low and
// high is returned if there are no results found.
func NodeDirs(kegpath string) (paths []string, low, high int) {
	low, high = -1, -1
	entries, err := os.ReadDir(kegpath)
	if err != nil {
		return
	}
	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() || !DirIsNode(name) {
			continue
		}
		val, _ := strconv.Atoi(filepath.Base(name))
		if low < 0 {
			low = val
		}
		if high < 0 {
			high = val
		}
		if val > high {
			high = val
		}
		if val < low {
			low = val
		}
		if abs, err := filepath.Abs(filepath.Join(kegpath, name)); err == nil {
			paths = append(paths, abs)
		}
	}
	return
}

func keys[M ~map[K]T, T any, K comparable](in M) []K {
	seen := map[K]bool{}
	for k, _ := range in {
		seen[k] = true
	}
	keys := []K{}
	for k, _ := range seen {
		keys = append(keys, k)
	}
	return keys
}
