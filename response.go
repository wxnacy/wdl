package wdl

import (
    "net/http"
    "fmt"
    "sort"
    "time"
    "strings"
)

type Response struct {
    Size int            // 总大小
    CompleteSize int    // 完成的大小
    FileName string     // 文件名
    Done chan bool      // 是否完成
    Res *http.Response  // http 请求的结果
    progressPrintPosition   int     // 输出位置
    StartTime time.Time    // 开始时间
    EndTime time.Time    // 结束时间
    Error chan error    // 错误
}

func (this *Response) Progress() float64{
    result := float64(this.CompleteSize) / float64(this.Size)
    return result
}

func (this *Response) IsComplete() bool {
    select {
        case <- this.Done: {
            return true
        }
        default: {
            return false
        }
    }
}

func (this *Response) PrintHeader() {
    res := this.Res
    status := Cyan(res.Status)
    if res.StatusCode != 200 {
        status = Red(res.Status)
    }
    fmt.Printf("%s %s\n", Blue(res.Proto), status)
    header := res.Header
    headerKeys := make([]string, 0, len(header))
    for k := range header {
        headerKeys = append(headerKeys, k)
    }
    sort.Strings(headerKeys)
    for _, k := range headerKeys {
        fmt.Printf("%s: %s\n", k, Cyan(header.Get(k)))
    }
    fmt.Println("")
    fmt.Println("")
}

func (this *Response) PrintData() {
    fmt.Println("data")
}

func (this *Response) Wait() {
    t := time.NewTicker(200 * time.Millisecond)
    defer t.Stop()

    Loop:
    for true {
        select {
            case <- t.C: {
                this.PrintProgress()
            }
            case <- this.Error: {
                e := <-this.Error
                fmt.Println(e)
                break Loop
            }
            case <- this.Done: {
                this.PrintProgress()
                break Loop
            }
        }
    }
}

// 打印下载过程
func (this *Response) PrintProgress() {
    if this.progressPrintPosition > 0 {
        fmt.Printf("\033[%dA\033[K", this.progressPrintPosition)
    }

    this.progressPrintPosition = 0
    fmt.Printf("%s \033[K\n", this.getProgressBar())
    this.progressPrintPosition++
}

// 获取下载进度条
func (this *Response) getProgressBar() string {
    barLen := 20                                            // 总长度
    progress := int(100 * this.Progress())

    progressBarLen := int(progress / int( 100 / barLen ))   // 进度条长度
    waitBarLen := barLen - progressBarLen                   // 等待条长度

    progressBar := strings.Repeat("=", progressBarLen)      // 进度条
    progressBar = Cyan(progressBar)

    waitBar := strings.Repeat("-", waitBarLen)              // 等待条
    waitBar = Yellow(waitBar)

    progressTime := time.Now().Sub(this.StartTime)

    return fmt.Sprintf(
        "%d / %d bytes [%s%s] %d%% %0.2fs",
        this.CompleteSize,
        this.Size,
        progressBar,
        waitBar,
        progress,
        progressTime.Seconds(),
    )
}
