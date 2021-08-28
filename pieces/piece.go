package pieces

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Piece struct {
	shape [][]bool
}

// NewFromBytes creates a new piece from bytes that represent a piece. The format is newline delimited combinations of
// either periods or spaces. Periods represent the shape the piece has and spaces can be used to align it properly.
// Other characters result in an error.
func NewFromBytes(raw []byte) (Piece, error) {
	var parsed [][]bool

	scanner := bufio.NewScanner(bytes.NewReader(raw))

	for scanner.Scan() {
		rawLine := scanner.Bytes()

		parsedLine := make([]bool, len(rawLine))
		for i, b := range rawLine {
			switch b {
			case ' ':
				parsedLine[i] = false
			case '.':
				parsedLine[i] = true
			default:
				return Piece{}, fmt.Errorf("unexpected character %c", b)
			}
			parsedLine[i] = b != ' '
		}

		parsed = append(parsed, parsedLine)
	}

	maxWidth := -1
	for _, row := range parsed {
		if len(row) > maxWidth {
			maxWidth = len(row)
		}
	}

	for i := range parsed {
		for len(parsed[i]) != maxWidth {
			parsed[i] = append(parsed[i], false)
		}
	}

	return Piece{
		shape: parsed,
	}, nil
}

// NewFromSlice creates a new piece from a 2D-slice of booleans that represent a piece. The slice should be in row major order.
// That is, indexing shape is done like shape[row][column], so e.g. shape[3][1] refers to the fourth row and second column.
// shape must be rectangular or an error is returned.
func NewFromSlice(shape [][]bool) (Piece, error) {
	width := -1
	for i, row := range shape {
		if i == 0 {
			width = len(row)
			continue
		}

		if len(row) != width {
			return Piece{}, errors.New("piece must be rectangular")
		}
	}

	return Piece{
		shape: shape,
	}, nil
}

// NewFromBytes creates a new piece from bytes that represent many piece. Each piece can be delimited by either newlines
// or lines that start with // (ignoring whitespace in both cases). For each piece, parsing is done as in NewFromBytes.
func NewFromManyBytes(raw []byte) ([]Piece, error) {
	var result []Piece

	s := bufio.NewScanner(bytes.NewReader(raw))

	var curr bytes.Buffer

	addPiece := func() error {
		raw := curr.Bytes()
		if len(raw) == 0 {
			return nil
		}

		piece, err := NewFromBytes(raw)
		if err != nil {
			return err
		}

		result = append(result, piece)
		return nil
	}

	for s.Scan() {
		line := s.Bytes()
		trimmedLine := bytes.TrimSpace(line)

		isSeparator := bytes.HasPrefix(trimmedLine, []byte{'/', '/'}) || len(trimmedLine) == 0

		if isSeparator {
			err := addPiece()
			if err != nil {
				return nil, err
			}

			curr.Reset()
		} else {
			curr.Write(line)
			curr.WriteRune('\n')
		}
	}

	err := addPiece()
	if err != nil {
		return nil, err
	}

	return result, nil
}

// MustFromSlice is like NewFromSlice but panics if an error occurs.
func MustFromSlice(shape [][]bool) Piece {
	result, err := NewFromSlice(shape)
	if err != nil {
		panic(err)
	}

	return result
}

// MustFromBytes is like NewFromBytes but panics if an error occurs.
func MustFromBytes(raw []byte) Piece {
	result, err := NewFromBytes(raw)
	if err != nil {
		panic(err)
	}

	return result
}

// MustFromManyBytes is like NewFromManyBytes but panics if an error occurs.
func MustFromManyBytes(raw []byte) []Piece {
	result, err := NewFromManyBytes(raw)
	if err != nil {
		panic(err)
	}

	return result
}

// Shape returns the shape of the piece.
func (p Piece) Shape() [][]bool {
	// TODO: don't clone?
	return p.clone(len(p.shape), len(p.shape[0]), func(_, _, origRowIndex, origColIndex int) (int, int) {
		return origRowIndex, origColIndex
	}).shape
}

// String return a human readable 2D representation of the piece.
func (p Piece) String() string {
	var s strings.Builder

	for _, row := range p.shape {
		for _, value := range row {
			if value {
				s.WriteRune('.')
			} else {
				s.WriteRune(' ')
			}
		}
		s.WriteRune('\n')
	}

	return s.String()
}

// RotateCW returns a new Piece that is this shape rotated clockwise. That is, turned 90 degrees from left to right.
func (p Piece) RotateCW() Piece {
	return p.clone(len(p.shape[0]), len(p.shape), func(numRows, numCols, origRowIndex, origColIndex int) (int, int) {
		return origColIndex, numCols - origRowIndex - 1
	})
}

// RotateCCW returns a new Piece that is this shape rotated counter-clockwise. That is, turned 90 degrees from right to left.
func (p Piece) RotateCCW() Piece {
	return p.clone(len(p.shape[0]), len(p.shape), func(numRows, numCols, origRowIndex, origColIndex int) (int, int) {
		return numRows - origColIndex - 1, origRowIndex
	})
}

// FlipVertical returns a new Piece that is this shape mirrored along the vertical axis. That is, the bottom becomes the top and vice versa.
func (p Piece) FlipVertical() Piece {
	return p.clone(len(p.shape), len(p.shape[0]), func(numRows, numCols, origRowIndex, origColIndex int) (int, int) {
		return numRows - origRowIndex - 1, origColIndex
	})
}

// FlipHorizontal returns a new Piece that is this shape mirrored along the horizontal axis. That is, the left becomes the right and vice versa.
func (p Piece) FlipHorizontal() Piece {
	return p.clone(len(p.shape), len(p.shape[0]), func(numRows, numCols, origRowIndex, origColIndex int) (int, int) {
		return origRowIndex, numCols - origColIndex - 1
	})
}

// clone returns a new Piece from this one according to an indexing function that translates row and column indices in
// this Piece to destination row and column indices in the clone.
func (p Piece) clone(numRows, numCols int, newIndices func(numRows, numCols, origRowIndex, origColIndex int) (int, int)) Piece {
	newShape := make([][]bool, numRows)

	for i := range newShape {
		newShape[i] = make([]bool, numCols)
	}

	for i, row := range p.shape {
		for j, v := range row {
			row, col := newIndices(numRows, numCols, i, j)
			newShape[row][col] = v
		}
	}

	return Piece{
		shape: newShape,
	}
}
