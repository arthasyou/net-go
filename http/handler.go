package http

import (
	"io/ioutil"
	"net/http"

	grpc "github.com/luobin998877/go_net/grpc"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	bin, _ := ioutil.ReadAll(r.Body)
	_, _, reply := grpc.Cli.SendJSON(1, 1, r.URL.Path, bin, 5)
	w.Write(reply)
}
