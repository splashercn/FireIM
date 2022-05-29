package client

import "net"

type Client struct {
	conn *net.Conn
}

func NewClient(conn *net.Conn) Client {
	c := Client{conn: conn}
	return c
}
