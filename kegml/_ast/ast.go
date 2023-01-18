package ast

// Node represents the output of a PEG packrat parser. Each comes with
// its own copy of the []rune slice (R) originally submitted to the
// parser. Even though all nodes have their own copies of the []rune
// slice the underlying array is identical. This is because slices are
// an abstraction in Go and all slices created from the same original
// data always share the same memory.
//
// The inclusive beginning of the matching data is preserved (B) as is
// the non-inclusive end (E). And for rules that have extraneous runes
// that should not be captured the beginning (CB) and end (CE) of the
// capture runes is also preserved.  For any given rule there can only
// be one capture range.
//
// For non-leaf nodes child nodes under this node are assigned to Under.
type Node struct {
	R     []rune `json:"-"` // copy of slice abstraction only (not underlying array)
	B     int    // beginning of capture (inclusive)
	E     int    // ending of capture (non-inclusive)
	XB    int    // beginning inclusive
	XE    int    // ending non-inclusive
	Under []Node // child nodes
}

func (n Node) String() string {
	// TODO something fancy if nodes under it.
	//	byt, err :=
}
