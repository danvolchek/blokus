package pieces_test

import (
	"danvolchek.com/blokus/pieces"
	"testing"
)

const (
	numStandardPieces = 21
)

func TestStandard(t *testing.T) {
	if len(pieces.Standard) != numStandardPieces {
		t.Errorf("got %d pieces, expected %d pieces", len(pieces.Standard), numStandardPieces)
	}
}
