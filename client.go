package Crane_RPC

import "net"

type Client struct {
	conn net.Conn
}

// NewClient creates a new client
func NewClient(conn net.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

func (c *Client) Call(serviceName string, args []interface{})  {

}