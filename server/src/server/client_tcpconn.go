package server

import (
	"encoding/binary"
)

//recv a msg package, block
func (c *Client) readTcp() (byte, []byte, error) {
	id := byte(0)
	if err := binary.Read(c.connTcp, binary.BigEndian, &id); err != nil {
		return 0, nil, err
	}
	b := make([]byte, msgLen(id))
	if err := binary.Read(c.connTcp, binary.BigEndian, b); err != nil {
		return 0, nil, err
	}
	return id, b, nil
}
