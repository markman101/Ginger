package recvbuffer

type RecvBuffer struct {
	_readPos  int
	_writePos int
	_data     []byte
	_dataLen  int
}

func NewRBuffer(bufSize int) *RecvBuffer {
	buf := &RecvBuffer{
		_readPos:  0,
		_writePos: 0,
		_data:     make([]byte, bufSize),
		_dataLen:  0,
	}
	return buf
}

func (rBuf *RecvBuffer) WritePos() []byte {
	return rBuf._data[rBuf._writePos:]
}

func (rBuf *RecvBuffer) ReadPos() []byte {
	return rBuf._data[rBuf._readPos:rBuf._writePos]
}

func (rBuf *RecvBuffer) RemainLen() int {
	return rBuf._dataLen - rBuf._writePos
}

func (rBuf *RecvBuffer) WriteOffsetAdd(add int) {
	rBuf._writePos = rBuf._writePos + add
}

func (rBuf *RecvBuffer) ReadOffsetAdd(add int) {
	rBuf._readPos = rBuf._readPos + add
}

func (rBuf *RecvBuffer) Reset() {
	copy(rBuf._data, rBuf._data[rBuf._readPos:rBuf._writePos])
	rBuf._writePos = rBuf._writePos - rBuf._readPos
	rBuf._readPos = 0
}
