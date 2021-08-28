package main

import (
	"danvolchek.com/blokus/pieces"
	"fmt"
)

func main() {
	//for _, piece := range pieces.Standard {
	//	fmt.Println(piece)
	//}

	fmt.Println(pieces.Standard[14])
	fmt.Println(pieces.Standard[14].FlipHorizontal())
	fmt.Println(pieces.Standard[14].FlipHorizontal().FlipHorizontal())
}
