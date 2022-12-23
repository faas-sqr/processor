package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"os/exec"
	grpc_client "processor/pkg/grpc-client"
	"processor/pkg/grpc-client/proto/function"
	"time"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	//初始化函数A
	functionAInit("functionA")

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

	r.Get("/functionA", func(w http.ResponseWriter, r *http.Request) {

		var m map[string]string
		m = make(map[string]string)
		m["key"] = "value1"
		m["key2"] = "value2"

		req := function.Event{
			Heads: m,
			Body:  "hello",
		}

		client, ok := grpc_client.ClientMap["functionA"]

		if !ok {
			//内部错误
			fmt.Errorf("client is not found...")
		}

		for i := 0; i < 10; i++ {
			r, err := client.ProcessEvent(ctx, &req)
			if err != nil {
				log.Fatalf("could not ProcessEvent: %v", err)
			}
			log.Printf("Body is: %s", r.GetBody())
		}

		w.Write([]byte("hello processor..."))
	})

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

func functionAInit(functionName string) {

	//可执行的函数启动二进制目录
	cmd := exec.Command("functionA-start", "/home/function/"+functionName)
	//cmd := exec.Command("/home/william/serverless/faas-sqr/processor/dockerfile/bin/functionA-start", "/home/function/"+functionName)

	e := cmd.Start()
	CheckError(e)
	//初始化A客户端
	grpc_client.ClientInit("functionA", "3000")
}
