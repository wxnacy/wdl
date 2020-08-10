package wdl

import (
    "time"
    "fmt"
)

var inProgress = 0

func Download(url string) {
    // begin := time.Now().Unix()

    client := Client{IsMultiLine: false, IsDownload: true}
    req := NewRequest(url)
    res := client.Do(req, "")

    res.Wait()

    // end := time.Now().Unix()
    // fmt.Println("time ", end-begin)
}

func DownloadBatch(urls ...string) {

    client := Client{IsDownload: true}

    resList := make([]*Response, len(urls), len(urls))

    var req *Request
    var res *Response

    for _, d := range urls {
        req = NewRequest(d)
        res = client.Do(req, "")
        resList = append(resList, res)

    }

    t := time.NewTicker(200 * time.Millisecond)
    defer t.Stop()

    Loop:
    for true {
        select {
            case <- t.C: {
                printProgress(resList)
            }
            case <- DoneBatch(resList): {
                printProgress(resList)
                break Loop
            }
        }
    }

}

func printProgress(resList []*Response) {
    progressPrintPosition := 0
    for _, d := range resList {

        if progressPrintPosition > 0 {
            fmt.Printf("\033[%dA\033[K", progressPrintPosition)
        }

        // this.progressPrintPosition = 0
        fmt.Printf("%s \033[K\n", d.getProgressBar())
        progressPrintPosition++
    }
}

func DoneBatch(resList []*Response) <-chan bool{
    result := make(chan bool, 1)
    result <- true
    for _, d := range resList {
        if ! d.IsComplete() {
            result <- false
        }
    }
    return result
}
