package main

import (
	"context"
	"fmt"
	_ "github.com/heroku/x/hmetrics/onload"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	fmt.Println("start port:", port)

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go run(ctx, wg)

	<-ctx.Done()

	wg.Wait()

	fmt.Println("Done...")
}

func run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-ticker.C:
			fmt.Println("Hello!")
		case <-ctx.Done():
			return
		}

	}
}
