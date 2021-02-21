package socket

import (
	"net"
	"os"
	"strconv"

	"github.com/arthasyou/utility-go/logger"
	"go.uber.org/zap"
)

// StartTCP server
func StartTCP(port uint16) {
	address := "0.0.0.0:" + strconv.Itoa(int(port))
	l, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error("Error listening: ", zap.String("err", err.Error()))
		os.Exit(1)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		// fmt.Println(c.RemoteAddr().String())
		if err != nil {
			logger.Error("Error connectiong:", zap.String("err", err.Error()))
			return
		}
		acceptor := newTCPAcceptor(node, c)
		go acceptor.read()
		go acceptor.write()
	}
}
