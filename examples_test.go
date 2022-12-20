package keg_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/rwxrob/keg"
)

func ExampleNewNode() {

	n := keg.NewNode()
	now := time.Now().UTC()

	fmt.Printf("ID: %q %v\n", n.ID, n.Changed.Sub(now) < 2*time.Second)

	// Output:
	// ID: "" true

}

func ExampleNewNodeFromLine() {

	n := keg.NewNodeFromLine("2\t2022-12-19 11:40:01Z\tSome title\t0,2\n")

	fmt.Println(`ID:`, n.ID)
	fmt.Printf("IntID: %T %v\n", n.IntID(), n.IntID())
	fmt.Println(`Changed:`, n.Changed)
	fmt.Println(`Title:`, n.Title)
	fmt.Println(`Nodes:`, n.Nodes)

	// Output:
	// ID: 2
	// IntID: int 2
	// Changed: 2022-12-19 11:40:01 +0000 UTC
	// Title: Some title
	// Nodes: [0 2]

}

func ExampleNewNodeFromLine_empty() {

	n := keg.NewNodeFromLine("")

	fmt.Printf("ID: %q\n", n.ID)
	fmt.Printf("Changed: %q\n", n.Changed)
	fmt.Printf("Title: %q\n", n.Title)
	fmt.Printf("Nodes: %q\n", n.Nodes)

	// Output:
	// ID: ""
	// Changed: "0001-01-01 00:00:00 +0000 UTC"
	// Title: ""
	// Nodes: []

}

func ExampleNewNodeFromLine_too_Many() {

	n := keg.NewNodeFromLine("2\t2022-12-19 11:40:01Z\tSome title\t0,2\tblah\tblah\n")

	fmt.Printf("ID: %q\n", n.ID)
	fmt.Printf("Changed: %q\n", n.Changed)
	fmt.Printf("Title: %q\n", n.Title)
	fmt.Printf("Nodes: %q\n", n.Nodes)

	// Output:
	// ID: ""
	// Changed: "0001-01-01 00:00:00 +0000 UTC"
	// Title: ""
	// Nodes: []

}

func ExampleNewNodeFromLine_negative_Allowed_But_Wrong() {

	n := keg.NewNodeFromLine("-20")

	// note no validation for field itself
	fmt.Printf("ID: %q\n", n.ID)

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

func ExampleFetchIndex() {

	text := "2\t2022-12-19 11:40:01Z\tSome title\t0,2\n" +
		"30\t2022-12-21 12:40:01Z\tSome other title\t2\n"

	// simulate server with a kegdex file
	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, text)
		})
	svr := httptest.NewServer(handler)
	defer svr.Close()
	dex, err := keg.FetchIndex(svr.URL)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%q\n", dex.Nodes[1])
	log.Println(dex.URL)

	dex, err = keg.FetchIndex(`bogus`)
	if err != nil {
		fmt.Println(err)
	}

	// simulate bad url
	handler2 := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		})
	svr2 := httptest.NewServer(handler2)
	defer svr2.Close()

	dex, err = keg.FetchIndex(svr2.URL)
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// "30\t2022-12-21 12:40:01Z\tSome other title\t2"
	// Get "bogus/kegdex": unsupported protocol scheme ""
	// failed to fetch: 400 Bad Request

}

func ExampleNewIndex() {

	dex := keg.NewIndex()

	fmt.Printf("%q\n", dex.File)
	fmt.Printf("%q\n", dex.URL)
	fmt.Printf("%q\n", dex.Nodes)

	//Output:
	// ""
	// ""
	// []

}

func ExampleReadIndex() {

	dex, err := keg.ReadIndex(`testdata/samplekeg`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dex.File)
	fmt.Println(len(dex.Nodes))
	fmt.Printf("%q\n", dex.Nodes[1])
	fmt.Printf("%q\n", dex.Nodes[2])

	// Output:
	// testdata/samplekeg/kegdex
	// 13
	// "1\t2022-11-26 19:33:24Z\tSample content node\t1,2,34,23"
	// "2\t2022-11-17 20:37:57Z\tSome title for 2"

}

func ExampleReadIndex_fail() {

	dex, err := keg.ReadIndex(`testdata/samlekeg`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dex)

	// Output:
	// open testdata/samlekeg/kegdex: no such file or directory
	// <nil>

}

func ExampleIndex_SortByID() {

	n1 := keg.NewNode()
	n1.ID = "0"
	n2 := keg.NewNode()
	n2.ID = "5"
	n3 := keg.NewNode()
	n3.ID = "3"

	dex := keg.NewIndex()
	dex.Add(n1, n2, n3)

	dex.SortByID()
	fmt.Println(dex.Nodes[1].ID)

	// Output:
	// 3
}

func ExampleIndex_SortByChanges() {

	n1 := keg.NewNode()
	n1.Changed = time.Now().UTC()
	n1.ID = "1"

	n2 := keg.NewNode()
	n2.Changed = n1.Changed.Add(-2 * time.Second)
	n2.ID = "2"

	n3 := keg.NewNode()
	n3.Changed = n1.Changed.Add(2 * time.Second)
	n3.ID = "3"

	dex := keg.NewIndex()
	dex.Add(n1, n2, n3)

	dex.SortByChanges()
	fmt.Println(dex.Nodes[0].ID)

	// Output:
	// 3
}

func ExampleIndex_MapIDs() {

	n1 := keg.NewNode()
	n1.Changed = time.Now().UTC()
	n1.ID = "1"

	n2 := keg.NewNode()
	n2.Changed = n1.Changed.Add(-2 * time.Second)
	n2.ID = "2"

	n3 := keg.NewNode()
	n3.Changed = n1.Changed.Add(2 * time.Second)
	n3.ID = "3"

	dex := keg.NewIndex()
	dex.Add(n3, n2, n1)

	fmt.Println(dex.IDs)
	dex.MapIDs()
	fmt.Println(len(dex.IDs))
	fmt.Println(dex.IDs["3"].ID)

	// Output:
	// map[]
	// 3
	// 3
}

func ExampleIndex_MapTitles() {

	n1 := keg.NewNode()
	n1.Changed = time.Now().UTC()
	n1.ID = "1"
	n1.Title = "Title of one"

	n2 := keg.NewNode()
	n2.Changed = n1.Changed.Add(-2 * time.Second)
	n2.ID = "2"
	n2.Title = "Title of two"

	n3 := keg.NewNode()
	n3.Changed = n1.Changed.Add(2 * time.Second)
	n3.ID = "3"
	n3.Title = "Title of three"

	dex := keg.NewIndex()
	dex.Add(n3, n2, n1)

	fmt.Println(dex.Titles)
	dex.MapTitles()
	fmt.Println(len(dex.Titles))
	fmt.Println(dex.Titles["Title of three"].ID)

	// Output:
	// map[]
	// 3
	// 3
}

func ExampleIndex_MapIncludes() {

	n1 := keg.NewNode()
	n1.Changed = time.Now().UTC()
	n1.ID = "1"
	n1.Title = "Title of one"
	n1.Nodes = []string{`3`, `2`}

	n2 := keg.NewNode()
	n2.Changed = n1.Changed.Add(-2 * time.Second)
	n2.ID = "2"
	n2.Title = "Title of two"
	n2.Nodes = []string{`2`}

	n3 := keg.NewNode()
	n3.Changed = n1.Changed.Add(2 * time.Second)
	n3.ID = "3"
	n3.Title = "Title of three"

	dex := keg.NewIndex()
	dex.Add(n3, n2, n1)

	fmt.Println(dex.Includes)
	dex.MapIncludes()
	fmt.Println(len(dex.Includes))
	fmt.Println(len(dex.Includes["2"]))
	fmt.Println(dex.Includes["2"]["1"].ID)
	fmt.Println(dex.Includes["2"]["2"].ID)
	fmt.Println(dex.Includes["2"]["3"])

	// Output:
	// map[]
	// 3
	// 2
	// 1
	// 2
	// <nil>
}

func ExampleIndex_String() {

	// also MarshalText

	n1 := keg.NewNode()
	n1.Changed, _ = time.Parse(`2006 Jan 2`, `2023 May 4`)
	n1.ID = "1"
	n1.Title = "Title of one"
	n1.Nodes = []string{`3`, `2`}

	n2 := keg.NewNode()
	n2.Changed = n1.Changed.Add(-2 * time.Second)
	n2.ID = "2"
	n2.Title = "Title of two"
	n2.Nodes = []string{`2`}

	n3 := keg.NewNode()
	n3.Changed = n1.Changed.Add(2 * time.Second)
	n3.ID = "3"
	n3.Title = "Title of three"

	dex := keg.NewIndex()
	dex.Add(n3, n2, n1)

	fmt.Printf("%q", dex)

	// Output:
	// "3\t2023-05-04 00:00:02Z\tTitle of three\t\n2\t2023-05-03 23:59:58Z\tTitle of two\t2\n1\t2023-05-04 00:00:00Z\tTitle of one\t3,2\n"

}

func ExampleIndex_Validate() {

	n1 := keg.NewNode()
	n1.Changed, _ = time.Parse(`2006 Jan 2`, `2023 May 4`)
	n1.ID = "1"
	n1.Title = ""
	n1.Nodes = []string{`3`, ``}

	n2 := keg.NewNode()
	n2.Changed = n1.Changed.Add(-2 * time.Second)
	n2.ID = ""
	n2.Title = "Title of two"
	n2.Nodes = []string{`2`}

	n3 := keg.NewNode()
	n3.Changed = n1.Changed.Add(2 * time.Second)
	n3.ID = "3"
	n3.Title = "Title of three ddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"

	n4 := new(keg.Node)
	n4.ID = "4"
	n4.Title = "Title of four"

	dex := keg.NewIndex()
	dex.Add(n3, n2, n1, n4)
	errors := dex.Validate()

	for _, err := range errors {
		fmt.Println(err)
	}

	// Output:
	// Title is too long: 134
	// Node identifier must be positive integer
	// Node title is empty
	// Node identifier must be positive integer
	// Node date last changed is not set (zero value)

}

func ExampleDirIsNode() {

	fmt.Println(keg.DirIsNode(`testdata`))
	fmt.Println(keg.DirIsNode(`testdata/samplekeg/1`))

	// Output:
	// false
	// true
}

func ExampleNodeDirs() {

	paths, low, high := keg.NodeDirs(`testdata/samplekeg`)
	fmt.Println(len(paths))
	fmt.Println(low)
	fmt.Println(high)

	// Output:
	// 13
	// 0
	// 12
}
