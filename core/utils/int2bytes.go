package utils

import (
	"bytes"
	"encoding/binary"
)

// BigEndianIntToBytes 整形转换成字节（大端）
func BigEndianIntToBytes(n int) ([]byte, error) {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	err := binary.Write(bytesBuffer, binary.BigEndian, x)
	if err != nil {
		return nil, err
	}
	return bytesBuffer.Bytes(), err
}

// BigEndianBytesToInt 字节转换成整形（大端）
func BigEndianBytesToInt(b []byte) (int, error) {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	err := binary.Read(bytesBuffer, binary.BigEndian, &x)
	if err != nil {
		return 0, err
	}

	return int(x), err
}

// LittleEndianIntToBytes 整形转换成字节（小端）
func LittleEndianIntToBytes(n int) ([]byte, error) {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	err := binary.Write(bytesBuffer, binary.LittleEndian, x)
	if err != nil {
		return nil, err
	}
	return bytesBuffer.Bytes(), err
}

// LittleEndianBytesToInt 字节转换成整形（小端）
func LittleEndianBytesToInt(b []byte) (int, error) {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	err := binary.Read(bytesBuffer, binary.LittleEndian, &x)
	if err != nil {
		return 0, err
	}

	return int(x), err
}
