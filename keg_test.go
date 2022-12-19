package keg

import (
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

func TestNode_String(t *testing.T) {
	n := &node{
		id:    "20",
		title: "Some title",
		nodes: []string{"3", "5"},
	}
	if n.String() != "20\t0001-01-01 00:00:00Z\tSome title\t3,5" {
		t.Error(`node conversion to string failing`)
	}
}

func TestNode_UnmarshalText(t *testing.T) {
	n := new(node)

	n.UnmarshalText([]byte("20\t2023-01-01 00:00:00Z\tSome title\t3,5"))

	switch {

	case n.id != "20":
		t.Error(`failed to unmarshal id`)

	case n.changed.Format(`2006-01`) != `2023-01`:
		t.Error(`failed to unmarshal change time stamp`)

	case n.title != "Some title":
		t.Error(`failed to unmarshal title`)

	case n.nodes[0] != "3" || n.nodes[1] != "5":
		t.Error(`failed to unmarshal nodes`)

	}

}
