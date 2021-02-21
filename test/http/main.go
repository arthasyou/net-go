package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc "github.com/luobin998877/go_net/grpc"
	"github.com/luobin998877/go_net/http"
)

func main() {
	grpc.Connect("127.0.0.1:8500", "testS")
	http.Start(8080)
	waitExit()
}

func waitExit() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	sig := <-ch

	fmt.Println("Got a signal", sig)
	now := time.Now()

	fmt.Println("Server exited", time.Since(now))
}
