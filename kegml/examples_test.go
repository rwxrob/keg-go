package kegml_test

import (
	"fmt"

	"github.com/rwxrob/keg/kegml"
)

func ExampleParseTitle() {

	//rat.Trace++
	fmt.Println(kegml.ParseTitle(`# Some title`))

	// Output:
	// Some title
}
