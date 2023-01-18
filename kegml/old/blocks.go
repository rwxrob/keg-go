package kegml

import (
	"github.com/rwxrob/keg/scan"
	"github.com/rwxrob/keg/scan/ast"
)

// ParseBlocks parses a KEGML README.md document
// into it's major blocks of runes.
//
//            BLOCK     │  TOKEN   │        CONTAINS
//      ────────────────┼──────────┼──────────────────────────
//        Title         │ #        │ Inflect, Math, Code
//        Bulleted List │ *  -  +  │ All but Lede
//        Numbered List │ 1.       │ All but Lede
//        Include List  │ * [      │ Inflect, Math, Code
//        Footnotes     │ [^       │ All buf Lede
//        Fenced        │ ``` ~~~  │ Runes
//        Quote         │ >        │ All but URL, Link, Lede
//        Math          │ $$       │ Runes
//        Figure        │ ![       │ Inflect, Math, Code
//        Separator     │ ----     │ None
//        Table         │ |        │ `
//        Paragraph     │ None     │ All but URL
//
func ParseBlocks(s *scan.R) *ast.Node {
	node := new(ast.Node)

	title := ParseTitle(s)

	if title != nil {
		node.Append(title)
	}

	return node
}
