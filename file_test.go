package wdl

import (
    "testing"
    "fmt"
)

func TestSortPieces(t *testing.T) {
    pieces := []Piece{Piece{Index: 3}, Piece{Index: 0}, Piece{Index: 2}}
    f := &File{Pieces: pieces}
    ps := f.SortPieces()

    if ps[0].Index != 0 {
        t.Errorf("%d 排序不正确", ps[0].Index)
    }
}

func TestAddPieceFromBytes(t *testing.T) {

    f := &File{}
    f.AddPieceFromBytes(1, []byte("Hello World"))
    f.AddPieceFromBytes(2, []byte("Hello Wxnacy"))

    if string(f.Pieces[0].Body) != "Hello World" {
        // t.Errorf("%v", f.Pieces[0].Body)
        fmt.Println("")
    }
}
