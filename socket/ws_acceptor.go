package socket

import (
	grpc "github.com/arthasyou/net-go/grpc"
	"github.com/arthasyou/net-go/packet"
	"github.com/arthasyou/utility-go/counter"
	"github.com/gorilla/websocket"
)

type wsChan struct {
	opCode int
	cmd    uint32
	data   []byte
	isEnd  bool
}

// Acceptor of web socket acceptor
type wsAcceptor struct {
	id   uint32
	conn *websocket.Conn
	ch   chan wsChan
}

func newWsAcceptor(node string, conn *websocket.Conn) *wsAcceptor {
	acceptor := wsAcceptor{
		id:   uint32(counter.Up("socket")),
		conn: conn,
		ch:   make(chan wsChan, 256),
	}
	return &acceptor
}

func (acceptor *wsAcceptor) read() {
	defer close(acceptor.ch)
	for {
		opCode, message, err := acceptor.conn.ReadMessage()
		if err != nil {
			acceptor.closeReader()
			return
		}
		cmd, dataBin, err := packet.H.UnpackData(message)
		if err != nil {
			acceptor.closeReader()
			return
		}
		acceptor.SendMsg(opCode, cmd, dataBin)
	}
}

func (acceptor *wsAcceptor) write() {
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
				grpc.Cli.SendSocket(node, acceptor.id, acceptor.conn.RemoteAddr().String(),
					1, 1, message.cmd, message.data, 5)
			bin := packet.H.PackData(code, rCmd, reply)
			err := acceptor.conn.WriteMessage(message.opCode, bin)
			if err != nil {
				break
			}
		}
	}
}

// SendMsg to grpc server
func (acceptor *wsAcceptor) SendMsg(opCode int, cmd uint32, msg []byte) {
	if acceptor == nil {
		return
	}
	acceptor.ch <- wsChan{opCode: opCode, cmd: cmd, data: msg}
	return
}

func (acceptor *wsAcceptor) closeReader() {
	acceptor.ch <- wsChan{isEnd: true}
}
