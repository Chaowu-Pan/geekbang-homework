package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	group, ctx := errgroup.WithContext(context.Background())

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World!")
	})

	//应用服务
	serverApp := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}

	//监控服务
	serverDebug := &http.Server{
		Addr: "127.0.0.1:8081",
	}

	servers := []*http.Server{serverApp, serverDebug}

	//启动应用服务
	group.Go(func() error {
		defer fmt.Println("Close some resources when App shutdown")
		return serverApp.ListenAndServe()
	})

	//启动监控服务
	group.Go(func() error {
		defer fmt.Println("Close some resources when Debug shutdown")
		return serverDebug.ListenAndServe()
	})

	group.Go(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		select {
		case <-c:
			//接收到系统信号
			shutdownServers(servers)
		case <-ctx.Done():
			//有服务返回错误
			shutdownServers(servers)
		}
		log.Println("server gracefully shutdown")
		return nil
	})

	if err := group.Wait(); err != nil {
		fmt.Printf("error: %v", err)
		return
	}

}

//关闭所有服务
func shutdownServers(servers []*http.Server) {
	for _, server := range servers {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("server shutdown failed, err: %v\n", err)
		}
		cancel()
	}
}
