package packet

// Handler creates a tcp handler that will be used to handle message.
type Handler interface {
	UnpackData(bin []byte) (cmd uint32, message []byte, err error)
	PackData(code uint32, cmd uint32, message []byte) []byte
}

// H packet handler for pack and unpack data
var H Handler

// Register with the same name, the one registered last will take effect.
func Register(h Handler) {
	H = h
}
