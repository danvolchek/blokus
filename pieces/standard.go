package pieces

import (
	_ "embed"
)

// Standard holds the set of standard pieces in a Blokus game.
var Standard []Piece

//go:embed standard.txt
var rawData []byte

func init() {
	Standard = MustFromManyBytes(rawData)
}
