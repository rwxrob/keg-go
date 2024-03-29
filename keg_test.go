package keg

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExists(t *testing.T) {
	if !exists("testdata/file") {
		t.Error(`exists failed to see existing file`)
	}
}

func TestNotExists(t *testing.T) {
	if !notExists("testdata/nope") {
		t.Error(`notExists fails to confirm missing file`)
	}
}

func TestFetch(t *testing.T) {

	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "2\t2022-12-19 11:40:01Z\tSome title\t0,2\n")
		})
	svr := httptest.NewServer(handler)
	defer svr.Close()

	buf, err := fetch(svr.URL)
	if err != nil {
		t.Error(err)
	}

	if string(buf) != "2\t2022-12-19 11:40:01Z\tSome title\t0,2\n" {
		t.Error(`fetch failed to get body`)
	}

}

type stringer struct{}

func (s stringer) String() string { return "stringer" }

func ExampleStringify() {
	r := strings.NewReader(`reader`)
	stuff := []any{
		"some", []byte{'t', 'h', 'i', 'n', 'g'},
		1, 2.234, stringer{}, r, []rune{'f', 'o', 'o'},
	}
	for _, s := range stuff {
		fmt.Printf("%q ", stringify(s))
	}
	// Output:
	// "some" "thing" "1" "2.234" "stringer" "reader" "foo"

}

func ExampleKeys() {
	it := map[string]bool{
		`one`:   true,
		`two`:   true,
		`three`: true,
	}

	for _, v := range keys(it) {
		fmt.Println(v)
	}
	// Unordered Output:
	// one
	// two
	// three
}

func ExampleAssertID() {
	fmt.Println(assertID("-1"))
	// Output:
	// Node identifier must be positive integer
}
