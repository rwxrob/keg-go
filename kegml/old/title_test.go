package kegml

import (
	"fmt"

	"github.com/rwxrob/keg/scan"
)

func ExampleParseTitle() {
	s := scan.New(`# Some good title`)
	n := ParseTitle(s)
	fmt.Println(n)
	// Output:
	// {"T":1,"V":"Some good title"}
}
