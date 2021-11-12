package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

type timeout interface {
	Timeout() bool
}

func handle(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 2)
	w.Header().Add("foo", "bar")
	w.WriteHeader(200)
}

func main() {
	http.HandleFunc("/", handle)
	go http.ListenAndServe(":9090", nil)

	c := &http.Client{
		Transport: &http.Transport{
			Proxy: nil,
			DialContext: (&net.Dialer{
				Timeout:   time.Second * 1, // 连接超时
				Deadline:  time.Time{},
				KeepAlive: 0,
			}).DialContext,
			WriteBufferSize:   4 << 10,  // 4KB
			ReadBufferSize:    10 << 10, // 10KB
			ForceAttemptHTTP2: false,
		},
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Second * 1,
	}

	req, _ := http.NewRequest(http.MethodGet, "http://0.0.0.0:9090/test", nil)
	rsp, err := c.Do(req)
	if err != nil {
		var e *url.Error
		if errors.As(err, &e) {
			log.Printf("timeout: %t [%s]", e.Timeout(), e.Error())
		}
		return
	}
	fmt.Printf("%+v ", rsp.StatusCode)
	fmt.Println(rsp.Header)
}
