package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/valyala/fasthttp"
	"net/http"
	"os"
	"processor/pkg/container"
	"processor/pkg/signals"
	"syscall"
)

func main() {
	container.StartFunction(os.Args[1])

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/kill-function", func(w http.ResponseWriter, r *http.Request) {

		err := container.CmdReexec.Process.Signal(syscall.SIGQUIT)
		if err != nil {
			//fmt.Println(err)
		}
		signals.Sigs <- syscall.SIGQUIT
		w.Write([]byte("kill function ok..."))
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		//测试函数
		w.Write([]byte(postFunction()))
	})

	server := &http.Server{
		Addr:    ":5000",
		Handler: r,
	}

	go listenSignal(context.Background(), server)
	err := server.ListenAndServe()
	if err != nil {
		return
	}
}

func listenSignal(ctx context.Context, httpSrv *http.Server) {

	select {
	case <-signals.Sigs:
		fmt.Println("notify sigs httpserver...")
		err := httpSrv.Shutdown(ctx)
		if err != nil {
			return
		}
		fmt.Println("http shutdown")
	}
}

func postFunction() string {
	url := "http://127.0.0.1:3000"
	status, resp, err := fasthttp.Get(nil, url)
	if err != nil {
		fmt.Println("请求失败:", err.Error())
		return "error-1"
	}

	if status != fasthttp.StatusOK {
		fmt.Println("请求没有成功:", status)
		return "error-2"
	}
	return string(resp)

}
