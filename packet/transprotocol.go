package packet

import (
	"encoding/binary"
	log "github.com/donnie4w/go-logger/logger"
)

type TransProtocol interface {
	Decode_pack(data []byte) (*Packet, int)

	Encode_pack(data []byte) *Packet
}

type TransProtocolComm struct {
}

func (protocol *TransProtocolComm) Decode_pack(data []byte) (*Packet, int) {
	log.Info("Enter Decode_pack")
	decodeLen := 0
	pack := &Packet{}
	dataLen := len(data)
	offset := 0
	headLen := binary.Size(pack._header)
	//1.get packhader
	//= not cause out of memory so >= is needed
	for offset = 0; dataLen-offset >= headLen; offset++ {
		pack._header._ver = int16(binary.BigEndian.Uint16(data[offset : offset+2]))
		pack._header._seq = int32(binary.BigEndian.Uint32(data[offset+2 : offset+6]))
		pack._header._type = int16(binary.BigEndian.Uint16(data[offset+6 : offset+8]))
		pack._header._bodyLen = int32(binary.BigEndian.Uint32(data[offset+8 : offset+12]))
		//get pack header successful
		if int16(VER) == pack._header._ver && pack._header._bodyLen < int32(MAX_PACK_SIZE) {
			break
		}
	}
	//2.check is a complete Packet
	if dataLen-offset < headLen || dataLen-offset < int(pack._header._bodyLen) {
		//can't find VER in data Slice or is not a complete pack.
		log.Error("[PACK_NOT_COMPLETE],pack bodyLen is", pack._header._bodyLen)
		return nil, 0
	}

	//3.get a complete pack
	pack._data = make([]byte, pack._header._bodyLen)
	copy(pack._data, data[offset:])
	decodeLen = offset + headLen + int(pack._header._bodyLen)
	return pack, decodeLen
}
func (protocol *TransProtocolComm) Encode_pack(data []byte) *Packet {
	pack := &Packet{}
	pack._header._ver = int16(VER)
	pack._header._seq = 1
	pack._header._type = 1
	pack._header._bodyLen = int32(len(data))
	pack._data = data
	return pack
}
