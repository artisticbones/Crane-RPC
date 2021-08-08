package Crane_RPC

import (
	"errors"
	"net"
	"reflect"
)

type Client struct {
	conn net.Conn
}

// NewClient creates a new client
func NewClient(conn net.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

// Call 将函数原型转换成函数
func (c *Client) Call(serviceName string, fptr interface{}) {
	container := reflect.ValueOf(fptr).Elem()

	// 在这个函数中完成对错误的处理，向服务器端建立连接，接收服务器端数据等操作
	// 客户端只需要调用该函数即可
	f := func(req []reflect.Value) []reflect.Value {
		// 创建连接
		clientTrans := NewTransport(c.conn)

		// 注册错误处理机制
		errorHandler := func(err error) []reflect.Value {
			// 依次处理返回参数
			outArgs := make([]reflect.Value, container.Type().NumOut())
			for i := 0; i < len(outArgs)-1; i++ {
				// Zero 返回一个值，表示指定类型的零值。
				// 结果与 Value 结构体的零值不同，零值表示根本没有值。
				// 例如， Zero(TypeOf(42)) 返回一个类型为 Int 且值为 0 的值。返回的值既不可寻址也不可设置。
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
			// 将错误信息放置在包末尾
			outArgs[len(outArgs)-1] = reflect.ValueOf(&err).Elem()
			return outArgs
		}
		// 处理包请求参数
		inArgs := make([]interface{}, 0, len(req))
		for i := range req {
			// 将请求参数作为接口类型加入，方便后续处理
			inArgs = append(inArgs, req[i].Interface())
		}
		// send request to server
		err := clientTrans.Send(Data{
			ServiceName: serviceName,
			Args:        inArgs,
		})
		if err != nil {
			return errorHandler(err)
		}
		// receive response from server
		rsp, err := clientTrans.Receive()
		if err != nil {
			return errorHandler(err)
		}
		if rsp.Err != "" {
			return errorHandler(errors.New(rsp.Err))
		}
		if len(rsp.Args) == 0 {
			rsp.Args = make([]interface{}, container.Type().NumOut())
		}
		// 解包响应包的参数
		numOut := container.Type().NumOut()
		outArgs := make([]reflect.Value, numOut)
		for i := 0; i < numOut; i++ {
			// unpackage arguments (except error)
			if i != numOut-1 {
				if rsp.Args[i] == nil {
					// if argument is nil (gob will ignore "Zero" in transmission), set "Zero" value
					outArgs[i] = reflect.Zero(container.Type().Out(i))
				} else {
					outArgs[i] = reflect.ValueOf(rsp.Args[i])
				}
			} else {
				// unpackdge error argument
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
		}
		return outArgs
	}

	// 根据container的类型，以及 Call 函数需要的功能构建的函数，构建出函数并赋给 container
	container.Set(reflect.MakeFunc(container.Type(), f))
}
