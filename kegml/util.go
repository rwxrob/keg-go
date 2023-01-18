package kegml

import (
	"fmt"
	"io"
	"log"
)

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
