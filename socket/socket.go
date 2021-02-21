package socket

import (
	"github.com/luobin998877/go_utility/counter"
)

var node string

var mTCP = make(map[uint32]*tcpAcceptor)

var mWs = make(map[uint32]*wsAcceptor)

var mPlayer = make(map[int]*player)

func init() {
	counter.Init("socket")
}

func registerTCP(a *tcpAcceptor) {
	mTCP[a.id] = a
}

func ungeginsterTCP(id uint32) {
	delete(mTCP, id)
}

func registerWS(a *wsAcceptor) {
	mWs[a.id] = a
}

func unreginsterWS(id uint32) {
	delete(mWs, id)
}

type player struct {
	id         int
	socketID   uint32
	socketType string
}

// RegisterNode use to service find node to boradcast
func RegisterNode(name string) {
	node = name
}

// RegisterPlayer to maps
func RegisterPlayer(id int, socketID uint32) {
	st := "tcp"
	if mTCP[socketID] == nil {
		st = "ws"
	}
	p := player{
		id:         id,
		socketID:   socketID,
		socketType: st,
	}
	mPlayer[id] = &p
}

// UnRegisterPlayer from maps
func UnRegisterPlayer(id int) {
	delete(mPlayer, id)
}

// BroadCast to players
func BroadCast(playerIDs []int, cmd uint32, data []byte) {
	for i := 0; i < len(playerIDs); i++ {
		p := mPlayer[i]
		if p.socketType == "tcp" {
			mTCP[p.socketID].SendMsg(cmd, data)
		} else {
			mWs[p.socketID].SendMsg(2, cmd, data)
		}
	}
}
