package grpc

import (
	"github.com/luobin998877/go_grpc_with_consul/service"
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
