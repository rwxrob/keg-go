package keg

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	IsoTimeLayout = `2006-01-02 15:04:05Z`
	UndefinedID   = `-1`
	IsoTimeExpStr = `\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\dZ`
)

// Node describes a node as contained in a knowledge exchange graph (KEG)
// or Index.
//
// ID
//
// The integer identifier as a string usually corresponding to the name of the
// directory containing the node. Although this is a string (since that
// is how it is most often used) it must always be an integer.
//
// Note that an ID of zero is a valid node identifier called the "zero
// node" according to the KEG specification used for linking to content
// that is yet to be completed.
//
//
// Title
//
// A string containing the title of the node not to exceed 70 runes
// (maximum of 280 bytes).
//
// Nodes
//
// Nodes returns all content node ids that this node depends on to exist.
// In KEGML terms, these are the node include links contained in an
// include list block. By traversing these a full content path through
// a keg can be obtained from any node that aggregates others.
//
// Nodes is guaranteed to always return a slice even if empty.
//
// Note that it is perfectly acceptable and expected that
// implementations of Node cast these Node interface instances back into
// specific private struct implementations within package scope to
// restore functionality otherwise lost by the limited, exported methods of
// the interface.
//
// String
//
// The fmt.Stringer interface is not required, but strongly recommended.
// It should call MarshalText.
//
// MarshalText
//
// The encoding.TextMarshaler interface is not required, but strongly
// recommended. See NewNodeFromLine for formatting of the line.
//
// UnmarshalText
//
// The encoding.TextUnmarshaler interface is not required, but strongly
// recommended. NewNodeFromLine should actually call UnmarshalText in
// most cases.
//
type Node struct {
	ID      string
	Title   string
	Changed time.Time
	Nodes   []string
}

// IntID converts the ID into a proper integer (usually using
// strconv.Atoi) or panics. Nodes containing non-integer IDs must
// produce fatal errors.
//
func (n Node) IntID() int {
	v, err := strconv.Atoi(n.ID)
	if err != nil {
		panic(err)
	}
	return v
}

// NewNode returns a pointer to a new Node struct with Changed time to
// now and initializing the Nodes with an empty slice.
func NewNode() *Node {
	n := new(Node)
	n.Changed = time.Now().UTC()
	n.Nodes = []string{}
	return n
}

// NewNodeFromLine takes a line of tab-delimited text and returns a
// new Node reference, which is always returned even
// if not all fields have been parsed.
//
// Fields are expected to be in the following order:
//
//     1. ID       - 0 or positive integer
//     2. Changed  - 2006-01-02 15:04:05Z
//     3. Title    - 70 runes maximum
//     4. Nodes    - integer strings (like ID) separated by commas
//
// Lines are passed through strings.TrimSpace and must be delimited with
// a single tab.
//
// Changed must be an time that matches IsoTimeLayout (2006-01-02 15:04:05Z).
//
// Nodes are expected to be positive integer strings (including 0) separated by
// a single comma with no spaces.
//
// Note that to remain performant no validation of the fields is done.
// See Index.Validate.
//
func NewNodeFromLine[T string | []byte | []rune](line T) *Node {
	n := new(Node)
	n.UnmarshalText([]byte(string(line)))
	return n
}

// UnmarshalText takes a line of tab-delimited text and unmarshals it.
// No error is ever returned as unmarshaling is on a best attempt basis.
// See NewNodeFromLine for details.
func (n *Node) UnmarshalText(text []byte) error {

	f := strings.Split(string(strings.TrimSpace(string(text))), "\t")

	length := len(f)
	if length > 4 {
		f = f[0:4]
	}

	switch length {

	case 4: // nodes
		n.Nodes = strings.Split(f[3], ",")
		fallthrough

	case 3: // title
		n.Title = f[2]
		fallthrough

	case 2: // changed
		n.Changed, _ = time.Parse(IsoTimeLayout, f[1])
		fallthrough

	case 1: // id
		n.ID = f[0]
	}

	return nil
}

func (n Node) MarshalText() ([]byte, error) {
	f := []string{n.ID, n.Changed.Format(IsoTimeLayout), n.Title}
	if n.Nodes != nil {
		f = append(f, strings.Join(n.Nodes, ","))
	}
	return []byte(strings.Join(f, "\t")), nil
}

func (n Node) String() string { b, _ := n.MarshalText(); return string(b) }

func assertID(id string) error {
	idn, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	if idn < 0 {
		return fmt.Errorf(_InvalidNodeID)
	}
	return nil
}

// Validate returns one error for every one of the following possible
// failed assertions:
//
//     * ID must be 0 or positive integer string
//     * Title must not be empty
//     * Nodes must all be valid IDs
//     * Changed must not be time.ZeroValue
//
func (n Node) Validate() []error {
	errors := make([]error, 0)

	if n.Title == "" {
		errors = append(errors, fmt.Errorf(_EmptyTitle))
	}

	if err := assertID(n.ID); err != nil {
		errors = append(errors, err)
	}

	for _, v := range n.Nodes {
		if err := assertID(v); err != nil {
			errors = append(errors, err)
		}
	}

	if n.Changed.IsZero() {
		errors = append(errors, fmt.Errorf(_ChangedIsZero))
	}

	return errors
}
