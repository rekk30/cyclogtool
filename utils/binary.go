package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

// ReadInt16 read signed 16 bit integer
func ReadInt16(data []byte) (ret int16) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

// ReadInt32 read signed 32 bit integer
func ReadInt32(data []byte) (ret int32) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

// ReadInt64 read signed 64 bit integer
func ReadInt64(data []byte) (ret int64) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

// ReadUInt16 read unsigned 16 bit integer
func ReadUInt16(data []byte) (ret uint16) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

// ReadUInt32 read unsigned 32 bit integer
func ReadUInt32(data []byte) (ret uint32) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

// ReadUInt64 read unsigned 64 bit integer
func ReadUInt64(data []byte) (ret uint64) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

// BinToUStr unsigned
func BinToUStr(data []byte, pos int, size int) (str string, result bool) {
	result = false
	str = ""
	if len(data) < pos+size {
		return
	}

	switch size {
	case 1:
		str = strconv.Itoa(int(data[pos]))
	case 2:
		str = strconv.Itoa(int(ReadUInt16(data[pos:])))
	case 4:
		str = strconv.Itoa(int(ReadUInt32(data[pos:])))
	case 8:
		str = strconv.Itoa(int(ReadUInt64(data[pos:])))
	default:
		fmt.Printf("Unsigned %d bytes param not implemented\n", size)
	}

	result = true
	return
}

// BinToStr signed
func BinToStr(data []byte, pos int, size int) (str string, result bool) {
	result = false
	str = ""
	if len(data) < pos+size {
		return
	}

	switch size {
	case 1:
		str = strconv.Itoa(int(data[pos]))
	case 2:
		str = strconv.Itoa(int(ReadInt16(data[pos:])))
	case 4:
		str = strconv.Itoa(int(ReadInt32(data[pos:])))
	case 8:
		str = strconv.Itoa(int(ReadInt64(data[pos:])))
	default:
		fmt.Printf("Signed %d bytes param not implemented\n", size)
	}

	result = true
	return
}
