package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Type byte

const (
	SimpleString Type = '+'
	BulkString   Type = '$'
	Array        Type = '*'
)

type Value struct {
	typ   Type
	bytes []byte
	array []Value
}

func (v Value) String() string {
	if v.typ == SimpleString || v.typ == BulkString {
		return string(v.bytes)
	}
	return ""
}

func (v Value) Array() []Value {
	if v.typ == Array {
		return v.array
	}
	return []Value{}
}

func DecodeRESP(byteStream *bufio.Reader) (Value, error) {
	dataTypeByte, err := byteStream.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch string(dataTypeByte) {
	case "+":
		return SimpleStringDecode(byteStream)
	case "$":
		return BulkStringDecode(byteStream)
	case "*":
		return ArrayDecode(byteStream)
	}

	return Value{}, fmt.Errorf("Invalid RESP data type byte: %s", string(dataTypeByte))
}

func SimpleStringDecode(byteStream *bufio.Reader) (Value, error) {
	readBytes, err := readUntilCRLF(byteStream)
	if err != nil {
		return Value{}, err
	}

	return Value{
		typ:   SimpleString,
		bytes: readBytes,
	}, nil
}

func BulkStringDecode(byteStream *bufio.Reader) (Value, error) {
	readBytesForCount, err := readUntilCRLF(byteStream)
	if err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string length: %s", err)
	}
	count, err := strconv.Atoi(string(readBytesForCount))
	if err != nil {
		return Value{}, fmt.Errorf("failed to parse bulk string length: %s", err)
	}

	readBytes := make([]byte, count+2)
	if _, err := io.ReadFull(byteStream, readBytes); err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string contents: %s", err)
	}

	return Value{
		typ:   BulkString,
		bytes: readBytes[:count],
	}, nil
}

func ArrayDecode(byteStream *bufio.Reader) (Value, error) {
	readBytesForCount, err := readUntilCRLF(byteStream)
	if err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string length: %s", err)
	}
	count, err := strconv.Atoi(string(readBytesForCount))
	if err != nil {
		return Value{}, fmt.Errorf("failed to parse bulk string length: %s", err)
	}

	array := []Value{}

	for i := 1; i <= count; i++ {
		value, err := DecodeRESP(byteStream)
		if err != nil {
			return Value{}, err
		}

		array = append(array, value)
	}

	return Value{
		typ:   Array,
		array: array,
	}, nil
}

func readUntilCRLF(byteStream *bufio.Reader) ([]byte, error) {
	readBytes := []byte{}
	for {
		b, err := byteStream.ReadBytes('\n')
		if err != nil {
			return nil, err
		}
		readBytes = append(readBytes, b...)
		if len(readBytes) >= 2 && readBytes[len(readBytes)-2] == '\r' {
			break
		}
	}

	return readBytes[:len(readBytes)-2], nil
}
