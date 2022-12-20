package keg

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

const IndexFileName = `keg-index`

type Index struct {
	File  string
	Nodes []*Node
}

// NewIndex returns an initialized Index with an empty slice of Nodes.
// Use ReadIndex to read an existing index from a KEG directory. Use
// ParseIndex to parse a new index directly from memory. Use FetchIndex
// to fetch an index from a URL.
func NewIndex() *Index {
	dex := new(Index)
	dex.Nodes = []*Node{}
	return dex
}

// ReadIndex examines the kegpath indicated for a file matching
// IndexFileName and if found buffers the file and returns ParseIndex.
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

	return dex, nil
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
