package keg

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
