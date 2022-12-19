package keg_test

import (
	"fmt"
	"time"

	"github.com/rwxrob/keg"
)

func ExampleNewNode() {

	n := keg.NewNode()
	now := time.Now().UTC()

	fmt.Printf("ID: %q %v\n", n.ID(), n.Changed().Sub(now) < 2*time.Second)

	// Output:
	// ID: "-1" true

}

func ExampleNewNodeFromLine() {

	n := keg.NewNodeFromLine("2\t2022-12-19 11:40:01Z\tSome title\t0,2\n")

	fmt.Println(`ID:`, n.ID())
	fmt.Printf("IntID: %T %v\n", n.IntID(), n.IntID())
	fmt.Println(`Changed:`, n.Changed())
	fmt.Println(`Title:`, n.Title())
	fmt.Println(`Nodes:`, n.Nodes())

	// Output:
	// ID: 2
	// IntID: int 2
	// Changed: 2022-12-19 11:40:01 +0000 UTC
	// Title: Some title
	// Nodes: [0 2]

}

func ExampleNewNodeFromLine_empty() {

	n := keg.NewNodeFromLine("")

	fmt.Printf("ID: %q\n", n.ID())
	fmt.Printf("Changed: %q\n", n.Changed())
	fmt.Printf("Title: %q\n", n.Title())
	fmt.Printf("Nodes: %q\n", n.Nodes())

	// Output:
	// ID: ""
	// Changed: "0001-01-01 00:00:00 +0000 UTC"
	// Title: ""
	// Nodes: []

}

func ExampleNewNodeFromLine_too_Many() {

	n := keg.NewNodeFromLine("2\t2022-12-19 11:40:01Z\tSome title\t0,2\tblah\tblah\n")

	fmt.Printf("ID: %q\n", n.ID())
	fmt.Printf("Changed: %q\n", n.Changed())
	fmt.Printf("Title: %q\n", n.Title())
	fmt.Printf("Nodes: %q\n", n.Nodes())

	// Output:
	// ID: ""
	// Changed: "0001-01-01 00:00:00 +0000 UTC"
	// Title: ""
	// Nodes: []

}

func ExampleNewNodeFromLine_negative_Allowed_But_Wrong() {

	n := keg.NewNodeFromLine("-20")

	// note no validation for field itself
	fmt.Printf("ID: %q\n", n.ID())

	// but will panic if attempted with IntID
	fmt.Printf("IntID: %T %v\n", n.IntID(), n.IntID())

	// Output:
	// ID: "-20"
	// IntID: int -20

}

func ExampleNewNodeFromLine_bad_ID() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	// parses fine
	n := keg.NewNodeFromLine("twenty")
	fmt.Printf("%q\n", n)

	// but panics when IntID attempted
	n.IntID()

	// Output:
	// "twenty\t0001-01-01 00:00:00Z\t"
	// strconv.Atoi: parsing "twenty": invalid syntax

}
