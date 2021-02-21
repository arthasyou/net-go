package main

import (
	"bytes"
	"encoding/binary"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc "github.com/arthasyou/net-go/grpc"
	"github.com/arthasyou/net-go/packet"
	"github.com/arthasyou/net-go/socket"

	"github.com/arthasyou/utility-go/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	t()
	initLog()
	initModules()

	waitExit()
}

func initLog() {
	settleLog := viper.GetString("logs/sys.log")
	logLevel := viper.GetString("debug")
	logger.InitLog(settleLog, logLevel)
}

func initModules() {
	socket.StartWs(8080)
}

type handler struct{}

func t() {
	socket.RegisterNode("testnode")
	grpc.Connect("127.0.0.1:8500", "testS")
	packet.Register(&handler{})
}

func (h *handler) UnpackData(bin []byte) (cmd uint32, message []byte, err error) {
	buffer := bytes.NewBuffer([]byte{})
	size := len(bin)
	var a uint16 = 0
	var b uint32 = 2
	binary.Write(buffer, binary.BigEndian, size)
	binary.Write(buffer, binary.BigEndian, a)
	binary.Write(buffer, binary.BigEndian, b)
	binary.Write(buffer, binary.BigEndian, bin)
	return 34, buffer.Bytes(), nil
}

func (h *handler) PackData(code uint32, cmd uint32, data []byte) []byte {
	buffer := bytes.NewBuffer([]byte{})
	size := len(data) + 6
	binary.Write(buffer, binary.BigEndian, size)
	binary.Write(buffer, binary.BigEndian, code)
	binary.Write(buffer, binary.BigEndian, cmd)
	binary.Write(buffer, binary.BigEndian, data)
	return buffer.Bytes()
}

func waitExit() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	sig := <-ch

	logger.Info("Got a signal", zap.String("sig", sig.String()))
	now := time.Now()

	logger.Info("Server exited", zap.String("exit duration", time.Since(now).String()))
}
