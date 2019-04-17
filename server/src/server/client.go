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

	data *Data

	statusCh chan Status
	actionCh chan Action

	quitCh chan struct{}
	wait   *sync.WaitGroup
}

type Data struct {
	status *Status
	xF     float64
	yF     float64
}

type Status struct {
	speedX uint16 //uint16 from stream
	speedY uint16 //uint16 from stream
}

type Action struct {
	id    byte //byte from stream
	value byte //byte from stream
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		conn:     conn,
		ip:       conn.RemoteAddr().String(),
		data:     &Data{status: &Status{}},
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
		id, msg, err := c.recv()
		if err != nil {
			if err != io.EOF {
				LogError("Client[%s] recv err:%v", c.ip, err)
			}
			//TODO notify quit
			return
		}
		//decode
		if isStatus(id) {
			status := decodeStatus(msg)
			c.statusCh <- status
		} else {
			action := decodeAction(msg)
			action.id = id
			c.actionCh <- action
		}
	}
}

func (c *Client) Serve() {
	ticker := time.NewTicker(time.Second / UpdateHZ)
	for {
		select {
		case st := <-c.statusCh:
			//update status
			c.data.status.speedX = st.speedX
			c.data.status.speedY = st.speedY
		case ac := <-c.actionCh:
			//do action
			doAction(ac)
		case <-ticker.C:
			//do status
			doStatus(c.data)
		case <-c.quitCh:
			c.wait.Done()
			return
		}
	}
}
