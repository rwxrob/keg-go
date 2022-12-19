package libkeg

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func notExists(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

const IsoTimeLayout = `2006-01-02 15:04:05Z`
const IsoTimeExpStr = `\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\dZ`

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
// IntID
//
// Converts the ID into a proper integer (usually using strconv.Atoi) or
// panics. Nodes containing non-integer IDs must produce fatal errors.
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
type Node interface {
	ID() string
	IntID() int
	Changed() time.Time
	Title() string
	Nodes() []string
}

type node struct {
	id      string
	title   string
	changed time.Time
	nodes   []string
}

func (n node) ID() string         { return n.id }
func (n node) Title() string      { return n.title }
func (n node) Changed() time.Time { return n.changed }

func (n node) Nodes() []string {
	if n.nodes == nil {
		return []string{}
	}
	return n.nodes
}

func (n node) IntID() int {
	v, err := strconv.Atoi(n.id)
	if err != nil {
		panic(err)
	}
	return v
}

// NewNodeFromLine takes a line of tab-delimited text and returns a node
// from it. A pointer to a struct that implements the Node interface is
// always returned even if the struct fields are empty.
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
// Note that no validation of the fields is done at all.
//
func NewNodeFromLine[T string | []byte | []rune](line T) *node {
	n := new(node)

	f := strings.Split(string(strings.TrimSpace(string(line))), "\t")

	length := len(f)
	if length > 4 {
		f = f[0:4]
	}

	switch length {

	case 4: // nodes
		n.nodes = strings.Split(f[3], ",")
		fallthrough

	case 3: // title
		n.title = f[2]
		fallthrough

	case 2: // changed
		n.changed, _ = time.Parse(IsoTimeLayout, f[1])
		fallthrough

	case 1: // id
		n.id = f[0]
	}

	return n
}
