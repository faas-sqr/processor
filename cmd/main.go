package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"net/http"
	"os/exec"
)

func main() {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		//测试函数
		w.Write([]byte("hello processor..."))
	})

	cmd := exec.Command("functionA-start", "/home/function/functionA")

	e := cmd.Start()
	CheckError(e)

	cmd = exec.Command("functionB-start", "/home/function/functionB")

	e = cmd.Start()
	CheckError(e)

	//processor 服务
	server := &http.Server{
		Addr:    ":9000",
		Handler: r,
	}
	err := server.ListenAndServe()
	if err != nil {
		return
	}

}

func CheckError(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
