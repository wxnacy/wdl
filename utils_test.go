package wdl

import (
    "testing"
    "os"
    "strings"
)

var pwd, _ = os.Getwd()
var filename = pwd + "/test/data/" + "file.txt"

func TestRenameFilePath(t *testing.T) {
    name := RenameFilePath(filename)
    if name != strings.Replace(filename, ".txt", "(2).txt", -1) {
        t.Errorf("%s 名称不正确", name)
    }

    name1 := RenameFilePath(name)
    if name1 != name {
        t.Errorf("%s 名称不正确", name)
    }

    name2 := RenameFilePath(pwd + "/test/data/test.txt")
    if name2 != pwd + "/test/data/test.txt" {
        t.Errorf("%s 名称不正确", name)
    }

}

func TestGetFileSuffix(t *testing.T) {
    suffix := GetFileSuffix(filename)
    if suffix != ".txt" {
        t.Errorf("%s 后缀不对", suffix)
    }
}

func TestGetFileNameByPath(t *testing.T) {
    suffix := GetFileNameByPath(filename)
    if suffix != "file.txt" {
        t.Errorf("%s 文件名不正确", suffix)
    }
}
