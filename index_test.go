package keg

import (
	"testing"
)

func TestIndex_MarshalText(t *testing.T) {
	dex := new(Index)
	buf, err := dex.MarshalText()
	if string(buf) != `` {
		t.Error(`failed to MarshalText`)
	}
	if err != nil {
		t.Error(err)
	}

}
