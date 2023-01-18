package kegml

import (
	"path/filepath"
	"strings"

	"github.com/rwxrob/keg/scan"
	"github.com/rwxrob/keg/scan/ast"
)

// ScanTitle scans a valid title (optionally parsing it into buf).
// A valid title must match the following:
//
// * Begin with hashtag and single space
// * <=70 unicode.IsPrint runes
//
// If more than 70 runes are available returns only 70. Note that this
// means further validation is generally needed on a README.md file to ensure it complies with the 70 rune specification.
//
func ScanTitle(s *scan.R, buf *[]rune) bool {
	m := s.Mark()
	if !s.Scan() || s.Rune() != '#' {
		return s.Revert(m, Title)
	}
	if !s.Scan() || s.Rune() != ' ' {
		return s.Revert(m, Title)
	}
	var count int
	for s.Scan() {
		if count >= 70 {
			return s.Revert(m, Title)
		}
		r := s.Rune()
		if r == '\n' {
			if count > 0 {
				return true
			} else {
				return s.Revert(m, Title)
			}
		}
		if buf != nil {
			*buf = append(*buf, r)
		}
		count++
	}
	return true
}

func ParseTitle(s *scan.R) *ast.Node {
	buf := make([]rune, 0, 70)
	if !ScanTitle(s, &buf) {
		return nil
	}
	return &ast.Node{T: Title, V: string(buf)}
}

// ReadTitle reads a KEG node title from KEGML file.
func ReadTitle(path string) (string, error) {
	if !strings.HasSuffix(path, `README.md`) {
		path = filepath.Join(path, `README.md`)
	}
	if err := Scanner.Open(path); err != nil {
		return "", err
	}
	nd := ParseTitle(Scanner)
	if nd == nil {
		return "", Scanner
	}
	return nd.V, nil
}
