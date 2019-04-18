package main

import (
	"fmt"
	"net"
	"encoding/binary"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:1701")
	if err != nil {
		panic(err)
	}

	var id byte
	var x int16
	var y int16
	b := make([]byte, 5)
	for {
		//三个参数：1 x y
		fmt.Scanln(&id, &x, &y)
		b[0] = id
		binary.BigEndian.PutUint16(b[1:], uint16(x))
		binary.BigEndian.PutUint16(b[3:], uint16(y))
		_, err := conn.Write(b)
		if err != nil {
			panic(err)
		}
	}
}
