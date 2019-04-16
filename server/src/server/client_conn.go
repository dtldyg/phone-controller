package server

import (
	"encoding/binary"
)

//recv a msg package, block
func (c *Client) recv() ([]byte, error) {
	len := uint16(0)
	if err := binary.Read(c.conn, binary.BigEndian, &len); err != nil {
		return nil, err
	}
	b := make([]byte, len)
	if err := binary.Read(c.conn, binary.BigEndian, b); err != nil {
		return nil, err
	}
	return b, nil
}
