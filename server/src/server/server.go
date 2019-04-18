package server

import (
	"net"
)

//TODO accept & manage
type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

//block
func (s *Server) Start() {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:1701")
	if nil != err {
		panic(err)
	}
	l, err := net.ListenTCP("tcp", addr)
	if nil != err {
		panic(err)
	}
	LogInfo("Server listening at [%v]", addr)

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			panic(err)
		}
		c := NewClient(conn)
		LogInfo("Client[%s] online", c.ip)

		go func() {
			c.Start()
			LogInfo("Client[%s] offline", c.ip)
		}()
	}
}
