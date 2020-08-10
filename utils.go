package wdl

import (
    "fmt"
    "os"
    "strings"
)

func main() {
    fmt.Println("Hello World")
}

func RenameFilePath(filePath string) string {
    index := 0
    result := ""
    Loop:
    for true {
        suffix := GetFileSuffix(filePath)
        replaceSuffix := suffix
        if index > 0 {
            replaceSuffix = fmt.Sprintf("(%d)%s", index, suffix)
        }
        newPath := strings.Replace(filePath, suffix, replaceSuffix, -1)
        _, err := os.Stat(newPath)
        if err != nil {
            result = newPath
            break Loop
        }
        index++
    }

    return result
}

// 获取文件后缀
func GetFileSuffix(filePath string) string {
    splits := strings.Split(filePath, ".")
    return fmt.Sprintf(".%s", splits[len(splits) - 1])
}

// 通过地址获取文件名
func GetFileNameByPath(filePath string) string {
    splits := strings.Split(filePath, "/")
    return fmt.Sprintf("%s", splits[len(splits) - 1])
}
