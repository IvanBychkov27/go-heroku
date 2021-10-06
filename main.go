package main

import (
	"context"
	"fmt"
	_ "github.com/heroku/x/hmetrics/onload"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Application struct {
	server *http.Server
}

func New() *Application {
	app := &Application{}

	router := http.NewServeMux()
	router.HandleFunc("/", handler)

	app.server = &http.Server{}
	app.server.Handler = router

	return app
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	addr := "127.0.0.1:" + port

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("error create main listener: ", err)
	}
	defer ln.Close()

	app := New()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go app.run(cancel, wg, ln)

	<-ctx.Done()

	app.stopServer()
	wg.Wait()

	fmt.Println("Done...")
}

func (app *Application) run(cancel context.CancelFunc, wg *sync.WaitGroup, ln net.Listener) {
	defer wg.Done()
	defer cancel()

	fmt.Println("start main listen address", ln.Addr().String())

	err := app.server.Serve(ln)
	if err != nil && err.Error() != "http: Server closed" {
		fmt.Println("error serve ", err.Error())
	}

}

func (app *Application) stopServer() {
	fmt.Println("stop server...")
	err := app.server.Shutdown(context.Background())
	if err != nil {
		fmt.Println("error stop main server:", err.Error())
	}
}

func handler(resp http.ResponseWriter, req *http.Request) {
	t := time.Now().Format("02:04:06 15:04:05")
	res := "<html><body><H1> Time: " + t + " </H1></body></html>"
	resp.Write([]byte(res))
}
