/*

Package kegml contains the parser functions and data types to create the
abstract syntax trees for a KEGML README.md documents. KEGML is an
extremely simplified form of GitHub Flavored Markdown.

*/
package kegml

import (
	"log"
	"unicode"

	"github.com/rwxrob/rat"
	"github.com/rwxrob/rat/x"
)

/*
var NodeBlocksX = x.Rule{`NodeBlocks`, 1, x.Seq{
	x.Seq{"# ", x.Rule{`Title`, 2, x.Any{1, 70}}, x.},
}}
*/

// EndLine <- CR? LF
var EndLine = x.Seq{x.Opt{'\r'}, '\n'}

// EndBlock <- !. / (EndLine / !. ) / EndLine{2}
var EndBlock = x.One{x.End{}, x.Seq{EndLine, x.End{}}, x.Rep{2, EndLine}}

// Title <= uprint{1,70}
var Title = x.Name{`Title`, x.Mmx{1, 70, unicode.IsPrint}}

// TitleBlock <- '#' SP Title EndBlock
var TitleLine = x.Seq{'#', ' ', Title, EndBlock}

var titleG *rat.Grammar

func ParseTitle(in any) string {
	if titleG == nil {
		titleG = rat.Pack(TitleLine)
	}
	res := titleG.Scan(stringify(in))
	if res.X != nil {
		log.Println(res.X)
	}
	results := res.WithName(`Title`)
	if len(results) > 0 {
		return results[0].Text()
	}
	return ""
}

/*
type Parser struct {
	R []rune
	I int
}

// Open is shortcut for new(Parser).Open(path)
func Open(path string) *Parser { return new(Parser).Open(path) }

// Scan is shortcut for new(Parser).Buffer(in).
func Scan(in any) *Parser { return new(Parser).Buffer(in) }

// Buffer sets buffer (R) and resets current index (I) to 0. Input may
// be io.Reader, []byte, []rune, or string.  Buffer is typically called
// immediately after instantiating a new scanner
// (ex: s := new(Parser).Buffer(in)).  Returns self-reference.
func (p *Parser) Buffer(in any) *Parser {
	switch v := in.(type) {
	case string:
		p.R = []rune(v)
	case []byte:
		p.R = []rune(string(v))
	case []rune:
		p.R = v
	case io.Reader:
		p, _ := io.ReadAll(v)
		if b != nil {
			s.R = []rune(string(b))
		}
	}
	return s
}

// Open opens the file at path and loads it by passing to Buffer.
// Returns self-reference.
func (s *Parser) Open(path string) *Parser {
	f, _ := os.Open(path)
	if f == nil {
		return s
	}
	defer f.Close()
	return s.Buffer(f)
}
*/

/*
func EndLine(s []rune) int {

	switch {
	case s[0] == '\n':
		return 1
	case len(s) == 1 && s[0] == '\r':
		return 1
	case len(s) > 1 && (s[0] == '\r' && s[1] == '\n'):
		return 2
	}
	return 0

}

func EndBlock(s []rune) int {
	n := EndLine(s)

	// end of data
	if n == len(s) {
		return n
	}

	n2 := EndLine(s[n:])
	if n2 > 0 {
		return n + n2
	}

	return 0
}

func Title(s []rune)  {
	if len(s) < 3 || s[0] != '#' || s[1] != ' ' {
		return 0
	}
	return 0
}
*/
