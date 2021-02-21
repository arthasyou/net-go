package http

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Start http server
func Start(port uint16) {
	r := mux.NewRouter()
	r.Use(mux.CORSMethodMiddleware(r))
	// Routes consist of a path and a handler function.
	r.PathPrefix("/api").HandlerFunc(apiHandler).Methods("POST")
	// Bind to a port and pass our router in
	address := ":" + strconv.Itoa(int(port))
	go http.ListenAndServe(address, r)
}

// StartSsl http server with ssl
func StartSsl(port uint16, cert string, key string) {
	r := mux.NewRouter()
	r.Use(mux.CORSMethodMiddleware(r))
	// Routes consist of a path and a handler function.
	r.PathPrefix("/api").HandlerFunc(apiHandler).Methods("POST")
	// Bind to a port and pass our router in
	address := ":" + strconv.Itoa(int(port))
	go http.ListenAndServeTLS(address, cert, key, r)
}
