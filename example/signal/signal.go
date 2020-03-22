package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// 启动一个最简单的 http 服务器
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello world"))
	})
	server := &http.Server{Addr: ":8080", Handler: mux}
	go func() {
		fmt.Println(server.ListenAndServe())
	}()

	// 通过 Notify 方法来监听信号
	sigs := make(chan os.Signal)
	signal.Notify(sigs)

	// 监听信号
	c := <-sigs
	fmt.Println(c.String())

	// 关闭服务
	fmt.Println(server.Close())
	time.Sleep(20 * time.Second)
	fmt.Println("退出程序")
}
