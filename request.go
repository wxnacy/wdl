package wdl

import (
    "fmt"
    "net/http"
    "strconv"
)

type Request struct {
    Method string
    Url string
    CanRestore bool
    Size int
    Req *http.Request
}

func (this *Request) GetHeader() *Response {
    client := http.Client{}
    req, err:= http.NewRequest("HEAD", this.Url, nil)
    res, err := client.Do(req)
    fmt.Println(err)

    length := res.Header.Get("Content-Length")
    this.Size, _ = strconv.Atoi(length)
    this.Req = req

    result := &Response{
        Res: res,
        Size: this.Size,

    }

    ranges := res.Header.Get("Accept-Ranges")
    if ranges == "bytes" {
        this.CanRestore = true
    }

    return result
}

func NewRequest(url string) *Request {
    req := &Request{Url: url}
    // req.Init()
    return req
}

