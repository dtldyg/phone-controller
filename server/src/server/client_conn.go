package server

import (
	"encoding/binary"
)

//recv a msg package, block
func (c *Client) recv() (byte, []byte, error) {
	id := byte(0)
	if err := binary.Read(c.conn, binary.BigEndian, &id); err != nil {
		return 0, nil, err
	}
	b := make([]byte, msgLen(id))
	if err := binary.Read(c.conn, binary.BigEndian, b); err != nil {
		return 0, nil, err
	}
	return id, b, nil
}
