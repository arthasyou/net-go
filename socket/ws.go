package socket

import (
	"flag"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options

// StartWs websocket
func StartWs(port uint16) {
	addr := "0.0.0.0:" + strconv.Itoa(int(port))
	http.HandleFunc("/", handler)
	go http.ListenAndServe(addr, nil)
}

// StartWss websocket with ssl
func StartWss(port uint16, cert string, key string) {
	addr := "0.0.0.0:" + strconv.Itoa(int(port))
	http.HandleFunc("/", handler)
	go http.ListenAndServeTLS(addr, cert, key, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	acceptor := newWsAcceptor(node, c)
	go acceptor.write()
	acceptor.read()
}
