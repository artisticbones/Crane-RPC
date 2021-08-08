package Crane_RPC

import (
	"encoding/binary"
	"io"
	"net"
)

type Transport struct {
	conn net.Conn
}

// NewTransport create a tranport
func NewTransport(conn net.Conn) *Transport {
	return &Transport{
		conn: conn,
	}
}

// Send data
func (t *Transport) Send(req Data) error {
	// 将请求数据编码程字节
	b, err := encode(req)
	if err != nil {
		return err
	}
	// 设置发送包
	buf := make([]byte, 4+len(b))
	// 设置报文头，包体长度
	binary.BigEndian.PutUint32(buf[:4], uint32(len(b)))
	copy(buf[4:], b)
	_, err = t.conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

// Receive data
func (t *Transport) Receive() (Data, error) {
	header := make([]byte, 4)
	_, err := io.ReadFull(t.conn, header)
	if err != nil {
		return Data{}, err
	}
	// 解析响应包头部
	dataLen := binary.BigEndian.Uint32(header)
	data := make([]byte, dataLen)
	_, err = io.ReadFull(t.conn, data)
	if err != nil {
		return Data{}, err
	}
	// 调用decode方法解析包
	rsp, err := decode(data)
	return rsp, err
}
