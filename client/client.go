package client

import (
	"bufio"
	"fmt"
	"net"
)

// Client is a super simple TCP Socket client
type Client struct {
	port int
}

func NewClient(port int) *Client {
	return &Client{
		port: port,
	}
}

func (cli *Client) SendRequest(message string) (response string, err error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", cli.port))
	if err != nil {
		return
	}
	defer conn.Close()

	if _, err = fmt.Fprintln(conn, message); err != nil {
		return
	}

	response, err = bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return
	}

	return
}
