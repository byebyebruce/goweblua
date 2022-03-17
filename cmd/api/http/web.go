// Package api_http接口
package api_http

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/byebyebruce/goweblua/executor"
)

var (
	//go:embed index.html
	html string
	Mgr  *executor.Manager
)

// HTTPHandleFunc web接口
func HTTPHandleFunc(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	if "reload" == r.Form.Get("cmd") {
		Mgr.RecreateNewWorker()
		w.Write([]byte("ok"))
		return
	}

	if r.Method == "GET" {
		t, err := template.New("index").Parse(html)
		if err != nil {
			w.Write([]byte("error"))
			return
		}
		t.Execute(w, nil)

	} else if r.Method == "POST" {
		defer func() {
			if r := recover(); nil != r {
				if e, ok := r.(error); ok {
					fmt.Println(e.Error())
				} else {
					fmt.Println(r)
				}
			}

		}()

		buf, err := ioutil.ReadAll(r.Body)

		if nil != err {
			w.Write([]byte(err.Error()))
			return
		}

		defer r.Body.Close()

		ctx1, cancel := context.WithTimeout(r.Context(), time.Second*time.Duration(1))
		defer cancel()

		var resultStr []byte

		select {
		case <-ctx1.Done():
			resultStr = []byte(ctx1.Err().Error())
		case ret := <-Mgr.Process(ctx1, buf):
			if len(ret) > 0 {
				resultStr = ret
			} else {
				resultStr = []byte("error: runtime empty")
			}

		}

		w.Write(resultStr)
	}

}
