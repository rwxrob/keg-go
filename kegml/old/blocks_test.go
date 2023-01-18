package kegml

import (
	"log"
	"testing"

	"github.com/rwxrob/keg/scan"
)

const readme = `# Here is some title

Just a paragraph
with soft line returns.

* [Some title](../3)
* [Another title](0)

`

func TestParseBlocks(t *testing.T) {
	ast := ParseBlocks(scan.New(readme))
	log.Print(ast)
}
