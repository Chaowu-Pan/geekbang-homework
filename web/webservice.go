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
	serverApp := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}

	serverDebug := &http.Server{
		Addr: "127.0.0.1:8081",
	}

	servers := []*http.Server{serverApp, serverDebug}

	group.Go(func() error {
		return serverApp.ListenAndServe()
	})

	group.Go(func() error {
		return serverDebug.ListenAndServe()
	})

	group.Go(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		select {
		case <-c:
			shutdownServers(servers)
		case <-ctx.Done():
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

func shutdownServers(servers []*http.Server) {
	for _, server := range servers {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("server shutdown failed, err: %v\n", err)
		}
		cancel()
	}
}
