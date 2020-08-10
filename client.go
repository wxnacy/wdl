package wdl

import (
    "fmt"
    "os"
    "strings"
    "time"
    "io/ioutil"
    "net/http"
    "errors"
)

type Client struct {
    Url string              // 访问地址
    PortSize int            // 分片下载的尺寸
    FileName string         // 下载文件名
    IsMultiLine bool        // 是否多线下载
    IsDownload bool         // 是否下载
}

var DefaultClient = &Client{}

func (this *Client) Do(req *Request, filename string) *Response {
    res := req.GetHeader()

    res.StartTime = time.Now()
    res.Done = make(chan bool, 1)
    res.Error = make(chan error, 1)
    res.FileName = filename
    this.Url = req.Url

    fmt.Printf("Connecting %s\n", req.Url)
    if this.IsMultiLine {
        this.MultiLineDownload(res)
    } else {
        this.SingleLineDownload(res)
    }

    return res
}


func (this *Client) SingleLineDownload(res *Response) {
    portSize, portNum, lastPort := this.GetPort(res.Size)
    filename := this.SaveFileName("")

    go func () {

        for i := 0; i< portNum; i++ {

            begin := i * portSize
            end := begin + portSize - 1
            if i + 1 == portNum {
                end = begin + lastPort - 1
            }

            ranges := fmt.Sprintf("%d-%d", begin, end)

            body, err := this.PortRequest(ranges)
            if err != nil {
                res.Error <- err
            }

            file, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
            _, err = file.Write(body)
            if err != nil {
                panic(err)
            }
            err = file.Close()
            if err != nil {
                panic(err)
            }

            res.CompleteSize = end + 1
            if res.CompleteSize == res.Size {
                res.Done <- true
                res.EndTime = time.Now()
            }
        }
    }()
}

func (this *Client) MultiLineDownload(res *Response) {
    portSize, portNum, lastPort := this.GetPort(res.Size)
    filename := this.SaveFileName("")

    for i := 0; i< portNum; i++ {

        begin := i * portSize
        end := begin + portSize - 1
        if i + 1 == portNum {
            end = begin + lastPort - 1
        }
        ranges := fmt.Sprintf("%d-%d", begin, end)
        savePath := fmt.Sprintf("%d.%s", i, filename)

        localfile := NewFile(savePath)

        go func() {
            body, _ := this.PortRequest(ranges)
            localfile.AddPieceFromBytes(i, body)

            if localfile.PieceNum() == portNum {
                localfile.Save()
                res.Done <- true
                res.EndTime = time.Now()
            }


            // file, _ := os.OpenFile(savePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
            // _, err := file.Write(body)
            // if err != nil {
                // fmt.Println(err)
                // log.Fatal(err)
            // }
            // err = file.Close()
            // if err != nil {
                // fmt.Println(err)
                // log.Fatal(err)
            // }

            // res.CompleteSize = end + 1
            // files, err := ioutil.ReadDir(".")
            // if err != nil {
                // fmt.Println(err)
                // log.Fatal(err)
            // }

            // dfs := make([]string, 0, len(files))
            // for _, file := range files {
                // if strings.HasSuffix(file.Name(), ".mp4") {
                    // dfs = append(dfs, file.Name())
                // }
            // }
            // fmt.Println(dfs)

            // if len(dfs) ==  portNum {

                // res.Done <- true
                // res.EndTime = time.Now()

                // sort.Strings(dfs)
                // fmt.Println(dfs)

                // tf, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

                // for _, d := range dfs {

                    // b, err := ioutil.ReadFile("./"+d)
                    // if err != nil {
                        // log.Fatal(err)
                    // }
                    // tf.Write(b)
                // }
                // tf.Close()

            // }

        }()

    }

}

func (this *Client) PortRequest(ranges string) ([]byte, error) {
    fmt.Println(ranges)
    client := http.Client{}
    req, err := http.NewRequest("GET", this.Url, nil)
    req.Header.Set("Range", fmt.Sprintf("bytes=%s", ranges))
    // req.Header.Set("Connection", "close")
    if err != nil {
        return nil, err
    }
    res, err := client.Do(req)
    fmt.Println(res.Status)
    if res.StatusCode != http.StatusPartialContent {
        return nil, errors.New(res.Status)
    }
    defer res.Body.Close()
    resp := &Response{Res: res}
    resp.PrintHeader()
    if err != nil {
        return nil, err
    }
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }
    return body, nil
}

func (this *Client) GetPort(size int) (int, int, int) {
    if this.PortSize == 0 {
        this.PortSize = 2 * 1024 * 1024
    }

    num := size / this.PortSize
    lastPort := size % this.PortSize
    if lastPort > 0 {
        num++
    }

    return this.PortSize, num, lastPort
}

func (this *Client) SaveFileName(filename string) string {
    urlPorts := strings.Split(this.Url, "/")
    saveName := urlPorts[len(urlPorts) - 1]
    if filename != "" {
        saveName = filename
    }
    dir, _ := os.Getwd()
    filePath := fmt.Sprintf("%s/%s", dir, saveName)
    filePath = RenameFilePath(filePath)
    fmt.Printf("Download file on %s\n", filePath)
    return filePath
}
