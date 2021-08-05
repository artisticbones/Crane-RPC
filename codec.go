package Crane_RPC

import (
	"bytes"
	"encoding/gob"
)

type Data struct {
	ServiceName string	// 服务名称
	Args []interface{}	// 传递的参数
	Err string			// socket 的错误
}

func encode(data Data) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decode(b []byte) (Data, error) {
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)
	var data Data
	if err := decoder.Decode(&data); err != nil {
		return Data{}, err
	}
	return data, nil
}