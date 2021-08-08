package Crane_RPC

import (
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
)

type Server struct {
	addr  string
	funcs map[string]reflect.Value
}

func NewServer(addr string) *Server {
	return &Server{
		addr:  addr,
		funcs: make(map[string]reflect.Value),
	}
}

// Register 通过名称注册一个方法
func (s *Server) Register(serviceName string, f interface{}) {
	// 如果 Map 中存在该名函数，直接返回
	if _, ok := s.funcs[serviceName]; ok {
		return
	}
	// 将函数映射到 Map 中
	s.funcs[serviceName] = reflect.ValueOf(f)
}

func (s *Server) Run() {
	// 监听本地址
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Printf("listen on %s err: %v\n", s.addr, err)
		return
	}
	// 循环监听来自客户端的请求
	for {
		// 接收来自客户端的请求
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept err: %v\n", err)
			continue
		}

		// 启动线程处理客户端的请求
		go func() {
			srvTransport := NewTransport(conn)

			for {
				// 从客户端读取请求包
				req, err := srvTransport.Receive()
				if err != nil {
					if err != io.EOF {
						log.Printf("read err: %v\n", err)
					}
					return
				}
				// 通过名称获取方法名
				f, ok := s.funcs[req.ServiceName]
				// 如果方法不存在
				if !ok {
					e := fmt.Sprintf("func %s does not exist", req.ServiceName)
					log.Println(e)
					if err = srvTransport.Send(Data{ServiceName: req.ServiceName, Err: e}); err != nil {
						log.Printf("tranport write err: %v\n", err)
					}
					continue
				}
				log.Printf("func %s is called\n", req.ServiceName)
				// 否则解包请求包
				inArgs := make([]reflect.Value, len(req.Args))
				for i := range req.Args {
					inArgs[i] = reflect.ValueOf(req.Args[i])
				}
				// 反射请求的方法
				out := f.Call(inArgs)
				// 构建响应包参数
				outArgs := make([]interface{}, len(out)-1)
				for i := 0; i < len(out)-1; i++ {
					outArgs[i] = out[i].Interface()
				}
				// 构建 Err 参数
				var e string
				if _, ok := out[len(out)-1].Interface().(error); !ok {
					e = ""
				} else {
					e = out[len(out)-1].Interface().(error).Error()
				}

				// 将构建好的响应包发送给客户端
				err = srvTransport.Send(Data{
					ServiceName: req.ServiceName,
					Args:        outArgs,
					Err:         e,
				})
				if err != nil {
					log.Printf("transport write err: %v\n", err)
				}
			}
		}()

	}
}
