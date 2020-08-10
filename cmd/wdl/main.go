package main

import (
    "flag"
    "github.com/wxnacy/wdl"
)

var isDownload bool
var isMultiLine bool

func InitArgs() {
    flag.BoolVar(&isDownload, "d", false, "Url is download")
    flag.BoolVar(&isMultiLine, "w", false, "Url is multiline download")
    flag.Parse()

}

func main() {

    InitArgs()


    urls := flag.Args()
    url := urls[0]

    // for i := 30; i < 38; i++ {
        // fmt.Println(whttp.SetColor("message", 0, 0, i ))
    // }

    wdl.Download(url)


}
