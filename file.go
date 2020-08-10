package wdl

import (
    "os"
    "io/ioutil"
    "sort"
)

type File struct {
    Size int64      // 文件大小
    Path string     // 地址
    Pieces []Piece
}

type Piece struct {
    Path string
    Body []byte     // 数据
    Index int
}

func (this *Piece) Length() int {
    return len(this.Body)
}

func (this *File) PieceNum() int {
    return len(this.Pieces)
}

func (this *File) SortPieces() []Piece {
    sort.Slice(this.Pieces, func (i, j int) bool {
        return this.Pieces[i].Index < this.Pieces[j].Index
    })
    return this.Pieces
}

func (this *File) AddPieceFromPath(filepath string) error{

    _, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    body, err := ioutil.ReadFile(filepath)
    if err != nil {
        return err
    }

    piece := Piece{Path: filepath, Body: body}

    this.Pieces = append(this.Pieces, piece)

    return nil
}

func (this *File) AddPieceFromBytes(index int, body []byte) error {

    piece := Piece{Index: index, Body: body}

    this.Pieces = append(this.Pieces, piece)

    return nil
}

func (this *File) Save() error {
    pieces := this.SortPieces()
    file, err := os.OpenFile(this.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    for _, d := range pieces {
        file.Write(d.Body)
    }
    file.Close()
    return nil
}

func NewFile(path string) *File {
    return &File{
        Path: path,
    }
}
