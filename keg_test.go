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
