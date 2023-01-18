package kegml

import (
	"github.com/rwxrob/keg/scan"
	"github.com/rwxrob/keg/scan/ast"
)

// ParseIncludeIDs returns a node containing all the
// identifiers parsed from the node include links of any include list
// block.
//
// Include list blocks must begin with a star (*) followed by a single
// space and then a left bracket and must follow an empty line.
//
//     * [Title to 3](../3)
//     * [Title to another](4)
//
// The integer is drawn from the markdown target and must either be
// a valid integer or an integer following exactly two periods and
// a single slash.
//
// Also see ParseBlocks.
//
// Parsing errors are ignored simply resulting in an empty string.
//
func ParseIncludes(s *scan.R) *ast.Node {
	// TODO
	return nil
}
