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
	addr, err := net.ResolveTCPAddr("tcp", ":1701")
	if nil != err {
		panic(err)
	}
	l, err := net.ListenTCP("tcp", addr)
	if nil != err {
		panic(err)
	}
	LogInfo("Server listening at [%s]", getInnerIP())

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

func getInnerIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, value := range addrs {
		if ipnet, ok := value.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return ""
}
