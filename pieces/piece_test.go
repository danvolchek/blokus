package pieces_test

import (
	"danvolchek.com/blokus/pieces"
	"embed"
	"reflect"
	"strconv"
	"testing"
)

func TestNewFromBytes(t *testing.T) {
	cases := []struct {
		in       string
		expected pieces.Piece
	}{
		{
			in:       ".",
			expected: pieces.MustFromSlice([][]bool{{true}}),
		},
		{
			in:       ".\n.",
			expected: pieces.MustFromSlice([][]bool{{true}, {true}}),
		},
		{
			in:       ". .\n  .",
			expected: pieces.MustFromSlice([][]bool{{true, false, true}, {false, false, true}}),
		},
		{
			in:       "   .\n....\n.",
			expected: pieces.MustFromSlice([][]bool{{false, false, false, true}, {true, true, true, true}, {true, false, false, false}}),
		},
	}

	for i, test := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			result := pieces.MustFromBytes([]byte(test.in))
			if !reflect.DeepEqual(result.Shape(), test.expected.Shape()) {
				t.Errorf("input:\n%s\ngot:\n%s\nexpected:\n%s", test.in, result, test.expected)
			}
		})
	}
}

func TestPiece_String(t *testing.T) {
	cases := []struct {
		in       pieces.Piece
		expected string
	}{
		{
			in:       pieces.MustFromSlice([][]bool{{true}}),
			expected: ".\n",
		},
		{
			in:       pieces.MustFromSlice([][]bool{{true}, {true}}),
			expected: ".\n.\n",
		},
		{
			in:       pieces.MustFromSlice([][]bool{{true, false, true}, {false, false, true}}),
			expected: ". .\n  .\n",
		},
		{
			in:       pieces.MustFromSlice([][]bool{{false, false, false, true}, {true, true, true, true}, {true, false, false, false}}),
			expected: "   .\n....\n.   \n",
		},
	}

	for i, test := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			result := test.in.String()
			if result != test.expected {
				t.Errorf("input:\n%s\ngot:\n%s\nexpected:\n%s", test.in, result, test.expected)
			}
		})
	}
}

//go:embed testdata/transforms
var transforms embed.FS

func TestPiece_RotateCW(t *testing.T) {
	testTransform(t, "testdata/transforms/rotate_cw.txt", pieces.Piece.RotateCW)
}

func TestPiece_RotateCCW(t *testing.T) {
	testTransform(t, "testdata/transforms/rotate_ccw.txt", pieces.Piece.RotateCCW)
}

func TestPiece_FlipHorizontal(t *testing.T) {
	testTransform(t, "testdata/transforms/flip_horizontal.txt", pieces.Piece.FlipHorizontal)
}

func TestPiece_FlipVertical(t *testing.T) {
	testTransform(t, "testdata/transforms/flip_vertical.txt", pieces.Piece.FlipVertical)
}

func testTransform(t *testing.T, sourceFile string, f func(p pieces.Piece) pieces.Piece) {
	data, err := transforms.ReadFile(sourceFile)
	if err != nil {
		t.Fatalf("couldn't read transforms testdata: %s", err)
	}

	expected := pieces.MustFromManyBytes(data)

	curr := expected[0]

	for i := 1; i < len(expected); i++ {
		next := f(curr)

		if !reflect.DeepEqual(next, expected[i]) {
			t.Fatalf("transform:\n%s\ngot:\n%s\nexpected:\n%s", curr, next, expected[i])
		}

		curr = next
	}
}
