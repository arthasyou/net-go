package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc "github.com/luobin998877/go_net/grpc"
	"github.com/luobin998877/go_net/packet"
	"github.com/luobin998877/go_net/socket"
	_ "github.com/luobin998877/go_net/socket"

	"github.com/luobin998877/go_utility/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	t()
	// TODO: 配置中心，动态配置
	// initConfig()
	initLog()
	// printVersion()

	initModules()

	waitExit()
}

func initLog() {
	settleLog := viper.GetString("logs/sys.log")
	logLevel := viper.GetString("debug")
	logger.InitLog(settleLog, logLevel)
}

func initModules() {
	socket.StartTCP(40001)
}

type handler struct{}

func t() {
	socket.RegisterNode("testnode")
	grpc.Connect("127.0.0.1:8500", "testS")
	packet.Register(&handler{})
}

func (h *handler) UnpackData(bin []byte) (cmd uint32, message []byte, err error) {
	// fmt.Println("length", len(bin))
	if len(bin) < 8 {
		err = errors.New("data size err")
		return
	}
	cmd = binary.BigEndian.Uint32(bin[0:])
	message = bin[3:]
	err = nil
	return
}

func (h *handler) PackData(code uint32, cmd uint32, data []byte) []byte {
	buffer := bytes.NewBuffer([]byte{})
	size := len(data) + 8
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
