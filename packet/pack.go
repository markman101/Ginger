package packet

import (
	"encoding/binary"
)

const (
	MAX_PACK_SIZE = 8 * 1024
)

const VER = 0x0501

type PacketHeader struct {
	_ver     int16
	_seq     int32
	_type    int16
	_bodyLen int32 //body len
}

type Packet struct {
	_header PacketHeader
	_data   []byte
}

func (pack *Packet) GetData() []byte {
	return pack._data
}
func (pack *Packet) Serialize() []byte {
	bodyLen := len(pack._data)
	headLen := binary.Size(pack._header)
	packLen := headLen + bodyLen
	buf := make([]byte, packLen)
	binary.BigEndian.PutUint16(buf[0:2], uint16(pack._header._ver))
	binary.BigEndian.PutUint32(buf[2:6], uint32(pack._header._seq))
	binary.BigEndian.PutUint16(buf[6:8], uint16(pack._header._type))
	binary.BigEndian.PutUint32(buf[8:12], uint32(pack._header._bodyLen))
	copy(buf[headLen:packLen], pack._data)
	return buf
}
