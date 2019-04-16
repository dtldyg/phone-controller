package server

import (
	"net"
	"time"
	"sync"
	"io"
)

type Client struct {
	conn net.Conn
	ip   string

	status *Status

	statusCh chan Status
	actionCh chan Action

	quitCh chan struct{}
	wait   *sync.WaitGroup
}

type Status struct {
	speedX float64 //int32 from stream
	speedY float64 //int32 from stream
}

type Action struct {
	id byte //byte from stream
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		conn:     conn,
		ip:       conn.RemoteAddr().String(),
		status:   &Status{0, 0},
		statusCh: make(chan Status, 256),
		actionCh: make(chan Action, 64),
		quitCh:   make(chan struct{}),
		wait:     &sync.WaitGroup{},
	}
}

//run client, block
func (c *Client) Start() {
	RunGo(c.wait, c.Recv)
	RunGo(c.wait, c.Serve)
	c.wait.Wait()
}

func (c *Client) Close() {
	c.conn.Close()
	close(c.quitCh)
}

func (c *Client) Recv() {
	for {
		msg, err := c.recv()
		if err != nil {
			if err != io.EOF {
				LogError("Client[%s] recv err:%v", c.ip, err)
			}
			//TODO notify quit
			return
		}
		//decode
	}
}

func (c *Client) Serve() {
	ticker := time.NewTicker(time.Second / UpdateHZ)
	for {
		select {
		case st := <-c.statusCh:
			//update status
			c.status.speedX = st.speedX
			c.status.speedY = st.speedY
		case ac := <-c.actionCh:
			//do action
		case <-ticker.C:
			//do status
			//TODO 根据当前status的速度，刷新鼠标
		case <-c.quitCh:
			c.wait.Done()
			return
		}
	}
}
