package keg

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const IndexFileName = `kegdex`

// Index contains the most commonly used information from a knowledge
// exchange graph in ways that are easily searchable. Indexes are almost
// always contained within a tab-delimited file for the most performant,
// universal persistence and transferability possible.
//
// The Index struct can itself be used as something of a database
// through the use of the following on-demand index maps:
//
//     * IDs
//     * Titles
//     * Includes
//
// Each has a corresponding Map* method to trigger its update. Otherwise,
// these fields remain nil.
//
// The special SortByID and SortByChanges methods change the order of
// the Nodes slice directly.
//
type Index struct {
	File     string                      // if from file system
	URL      string                      // if from network
	Nodes    []*Node                     // order differs (SortByID, SortByChanges)
	IDs      map[string]*Node            // after calling MapIDs
	Titles   map[string]*Node            // after calling MapTitles
	Includes map[string]map[string]*Node // after calling MapIncludes
}

// SortByID sorts the Nodes slice by increasing numeric value of ID.
func (dex *Index) SortByID() {
	sort.Slice(dex.Nodes, func(i, j int) bool {
		return dex.Nodes[i].ID < dex.Nodes[j].ID
	})
}

// SortByChanges sorts the Nodes slice by most recent change (Changed in
// reverse chronological order).
func (dex *Index) SortByChanges() {
	sort.Slice(dex.Nodes, func(i, j int) bool {
		return dex.Nodes[i].Changed.After(dex.Nodes[j].Changed)
	})
}

// MapIDs updates the internal IDs map creating a map index keyed to IDs
// pointing to their Node references. Before this is called the IDs map
// is nil. Note that there is no validation of the content and that any
// nodes with the same ID will be overwritten. Consider calling Validate
// before when needed.
func (dex *Index) MapIDs() {
	dex.IDs = make(map[string]*Node, len(dex.Nodes))
	for _, n := range dex.Nodes {
		dex.IDs[n.ID] = n
	}
}

// MapTitles updates the internal Titles map creating a map index keyed to Titles
// pointing to their Node references. Before this is called the Titles map
// is nil. Note that there is no validation of the content and that any
// nodes with the same Title will be overwritten. Consider calling Validate
// before when needed.
func (dex *Index) MapTitles() {
	dex.Titles = make(map[string]*Node, len(dex.Nodes))
	for _, n := range dex.Nodes {
		dex.Titles[n.Title] = n
	}
}

// MapIncludes updates the internal Includes map creating a map index
// keyed to each node in the Nodes field providing a way to quickly
// lookup every node that includes a specific node for easy dependency
// tracking, visualization, and such. Before MapIncludes is called the
// Includes map is nil. Note that there is no validation that any of the
// nodes added to the Include map actually exist. Consider calling
// Validate before when needed.
func (dex *Index) MapIncludes() {
	dex.Includes = make(map[string]map[string]*Node, len(dex.Nodes))
	for _, n := range dex.Nodes {
		dex.Includes[n.ID] = map[string]*Node{}
		for _, in := range n.Nodes {
			dex.Includes[in][n.ID] = n
		}
	}
}

// MarshalText fulfills the encoding.TextMarshaler interface by
// returning the same tab-delimited text expected in any index file. An
// error is never returned and a byte slice, even if length of zero, is
// always returned.
func (dex Index) MarshalText() ([]byte, error) {
	var str string
	if dex.Nodes == nil {
		return []byte{}, nil
	}
	for _, n := range dex.Nodes {
		str += strings.Join([]string{
			n.ID, n.Changed.Format(IsoTimeLayout), n.Title, strings.Join(n.Nodes, ","),
		}, "\t") + "\n"
	}
	return []byte(str), nil
}

/*
// UnmarshalText fulfills the encoding.TextUnmarshaler interface while
// preserving references to existing values. UnmarshalText preserves
// referential integrity. Existing Nodes will have their fields updated
// if detected during unmarshaling. New Nodes will be added. Any non-nil
// map will be updated as well (IDs, Titles, Includes).
func (dex *Index) UnmarshalText(buf []byte) error {
	s := bufio.NewScanner(strings.NewReader(string(buf)))
	dex.MapIDs()
	for s.Scan() {
		node := s.Text()
		// TODO
	}
	return nil
}
*/

// String fulfills the fmt.Stringer interface by returning the same
// tab-delimited text expected in any index file. See MarshalText.
func (dex Index) String() string { b, _ := dex.MarshalText(); return string(b) }

// NewIndex returns an initialized Index with an empty slice of Nodes.
// Use ReadIndex to read an existing index from a KEG directory. Use
// ParseIndex to parse a new index directly from memory. Use FetchIndex
// to fetch an index from a KEG URL. Use ScanIndex to create a new index
// from a KEG directory containing KEG content node directories.
func NewIndex() *Index {
	dex := new(Index)
	dex.Nodes = []*Node{}
	return dex
}

// ReadIndex examines the kegpath indicated for a file matching
// IndexFileName and if found buffers the file and returns ParseIndex.
// Note that this does not trigger a scan and returns an error if there
// is no index file found.
func ReadIndex(kegpath string) (*Index, error) {

	var dex *Index
	var err error
	var buf []byte

	file := filepath.Join(kegpath, IndexFileName)

	buf, err = os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	dex, err = ParseIndex(buf)
	dex.File = file

	return dex, err

}

// ParseIndex parses any of the following into a new Index:
//
// * string
// * []byte
// * []rune
// * io.Reader
//
// To remain performant, each line of input is parsed and loaded as is
// without any validation. See Index.Validate.
//
func ParseIndex(in any) (*Index, error) {
	dex := new(Index)

	s := bufio.NewScanner(strings.NewReader(stringify(in)))

	for line := 1; s.Scan(); line++ {
		node := NewNodeFromLine(s.Text())
		dex.Nodes = append(dex.Nodes, node)
	}

	return dex, nil
}

// FetchIndex fetches the data from the target URL and passes it to
// ParseIndex returning any error and always returning an index pointer.
// If the URL does not end with IndexFileName then it is added.
func FetchIndex(kegurl string) (*Index, error) {

	url := kegurl + `/` + IndexFileName

	var (
		dex *Index
		err error
		buf []byte
	)

	buf, err = fetch(url)
	if err != nil {
		return nil, err
	}

	dex, err = ParseIndex(buf)
	if err != nil {
		return nil, err
	}

	dex.URL = url

	return dex, nil
}

// ScanIndex takes the path to a keg directory and scans all the
// directories with node ids for names. Each content node directory is
// passed to ReadNode and the new node is appended to the
// Nodes slice of the Index. An Index is always returned even if empty.
func ScanIndex(kegpath string) (*Index, error) {
	dex := new(Index)
	var err error
	// TODO
	return dex, err
}

// Validate iterates over every node calling Validate on it and adding
// any returning error to the slice it returns. There is no limit on the
// number of errors. Returns an empty slice if no errors encountered.
func (dex *Index) Validate() []error {
	errors := []error{}
	for _, node := range dex.Nodes {
		errors = append(errors, node.Validate()...)
	}
	return errors
}

// Add is a convenience method to append nodes to the internal Nodes
// slice. Panics if dex.Nodes is nil.
func (dex *Index) Add(nodes ...*Node) { dex.Nodes = append(dex.Nodes, nodes...) }
