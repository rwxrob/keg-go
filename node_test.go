package keg

import (
	"testing"
)

func TestNode_String(t *testing.T) {
	n := &Node{
		ID:       "20",
		Title:    "Some title",
		Includes: []string{"3", "5"},
	}
	if n.String() != "20\t0001-01-01 00:00:00Z\tSome title\t3,5" {
		t.Error(`node conversion to string failing`)
	}
}

func TestNode_UnmarshalText(t *testing.T) {
	n := new(Node)

	n.UnmarshalText([]byte("20\t2023-01-01 00:00:00Z\tSome title\t3,5"))

	switch {

	case n.ID != "20":
		t.Error(`failed to unmarshal id`)

	case n.Changed.Format(`2006-01`) != `2023-01`:
		t.Error(`failed to unmarshal change time stamp`)

	case n.Title != "Some title":
		t.Error(`failed to unmarshal title`)

	case n.Includes[0] != "3" || n.Includes[1] != "5":
		t.Error(`failed to unmarshal nodes`)

	}

}
