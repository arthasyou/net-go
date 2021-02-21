package socket

import (
	"encoding/binary"
	"net"

	grpc "github.com/luobin998877/go_net/grpc"
	"github.com/luobin998877/go_net/packet"
	"github.com/luobin998877/go_utility/counter"
)

const (
	headLen = 2
)

type chanData struct {
	cmd   uint32
	data  []byte
	isEnd bool
}

// TCPAcceptor of tcp
type tcpAcceptor struct {
	id   uint32
	node string
	conn net.Conn
	ch   chan chanData
}

func newTCPAcceptor(node string, conn net.Conn) *tcpAcceptor {
	acceptor := tcpAcceptor{
		id:   uint32(counter.Up("socket")),
		node: node,
		conn: conn,
		ch:   make(chan chanData, 256),
	}
	registerTCP(&acceptor)
	return &acceptor
}

func (acceptor *tcpAcceptor) read() {
	defer close(acceptor.ch)

	for {
		size, err := handleHead(acceptor.conn)
		if err != nil {
			acceptor.closeReader()
			return
		}
		bin, err := handleData(acceptor.conn, size)
		if err != nil {
			acceptor.closeReader()
			return
		}
		cmd, dataBin, err := packet.H.UnpackData(bin)
		if err != nil {
			acceptor.closeReader()
			return
		}
		acceptor.SendMsg(cmd, dataBin)
	}
}

func handleHead(conn net.Conn) (uint16, error) {
	buffer := make([]byte, headLen)
	_, err := conn.Read(buffer)
	if err != nil {
		conn.Close()
		return 0, err
	}
	size := binary.BigEndian.Uint16(buffer[0:])
	return size, nil
}

func handleData(conn net.Conn, size uint16) ([]byte, error) {
	buffer := make([]byte, size)
	_, err := conn.Read(buffer)
	if err != nil {
		conn.Close()
		return nil, err
	}
	return buffer, nil
}

func (acceptor *tcpAcceptor) write() {
	for {
		select {
		case message, ok := <-acceptor.ch:
			if !ok {
				return
			}
			if message.isEnd {
				return
			}
			code, _, _, rCmd, reply :=
				grpc.Cli.SendSocket(acceptor.node, acceptor.id, acceptor.conn.RemoteAddr().String(),
					1, 1, message.cmd, message.data, 5)
			bin := packet.H.PackData(code, rCmd, reply)
			_, err := acceptor.conn.Write(bin)
			if err != nil {
				break
			}
		}
	}
}

// SendMsg to grpc server
func (acceptor *tcpAcceptor) SendMsg(cmd uint32, msg []byte) {
	if acceptor == nil {
		return
	}
	acceptor.ch <- chanData{cmd: cmd, data: msg, isEnd: false}
	return
}

func (acceptor *tcpAcceptor) closeReader() {
	acceptor.ch <- chanData{isEnd: true}
}
