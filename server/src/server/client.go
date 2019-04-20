package server

import (
	"net"
	"time"
	"sync"
)

type Client struct {
	connTcp net.Conn
	connUdp *net.UDPConn
	ip      string

	status *Status

	statusCh chan Status
	actionCh chan Action

	quitCh chan struct{}
	wait   *sync.WaitGroup
}

type Status struct {
	moveX int16 //int16 from stream
	moveY int16 //int16 from stream
}

type Action struct {
	id    byte //byte from stream
	value byte //byte from stream
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		connTcp:  conn,
		ip:       conn.RemoteAddr().String(),
		status:   &Status{},
		statusCh: make(chan Status, 256),
		actionCh: make(chan Action, 64),
		quitCh:   make(chan struct{}),
		wait:     &sync.WaitGroup{},
	}
}

//run client, block
func (c *Client) Start() {
	addr, err := net.ResolveUDPAddr("udp", ":1702")
	if nil != err {
		panic(err)
	}
	c.connUdp, err = net.ListenUDP("udp", addr)
	if nil != err {
		panic(err)
	}
	RunGo(c.wait, c.RecvTcp)
	RunGo(c.wait, c.RecvUdp)
	RunGo(c.wait, c.Serve)
	c.wait.Wait()
}

func (c *Client) RecvTcp() {
	for {
		id, msg, err := c.readTcp()
		if err != nil {
			close(c.quitCh)
			c.connUdp.Close()
			return
		}
		c.dispatch(id, msg)
	}
}

func (c *Client) RecvUdp() {
	for {
		id, msg, err := c.readUdp()
		if err != nil {
			return
		}
		c.dispatch(id, msg)
	}
}

func (c *Client) dispatch(id byte, msg []byte) {
	if isStatus(id) {
		status := decodeStatus(msg)
		c.statusCh <- status
	} else {
		action := decodeAction(msg)
		action.id = id
		c.actionCh <- action
	}
}

func (c *Client) Serve() {
	ticker := time.NewTicker(time.Second / UpdateHZ)
	for {
		select {
		case st := <-c.statusCh:
			//update status
			c.status.moveX += st.moveX
			c.status.moveY += st.moveY
		case ac := <-c.actionCh:
			doAction(ac)
		case <-ticker.C:
			if c.status.moveX != 0 || c.status.moveY != 0 {
				doStatus(c.status)
				c.status.moveX = 0
				c.status.moveY = 0
			}
		case <-c.quitCh:
			doQuit()
			return
		}
	}
}
