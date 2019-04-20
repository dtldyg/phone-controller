package server

//recv a msg package, block
func (c *Client) readUdp() (byte, []byte, error) {
	b := make([]byte, 5)
	n, _, err := c.connUdp.ReadFromUDP(b)
	if err != nil {
		return 0, nil, err
	}
	return b[0], b[1:n], nil
}
