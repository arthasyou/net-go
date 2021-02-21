package grpc

import (
	"github.com/arthasyou/grpc-consul-go/service"
)

// Cli grpc client
var Cli *service.Client

// Connect to grpc server
func Connect(consulAddr string, serviceName string) {
	c, err := service.CreateConnection(consulAddr, serviceName)
	if err != nil {
		return
	}
	Cli = c
}
