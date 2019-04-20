package server

import "encoding/binary"

func msgLen(id byte) int {
	switch id {
	case 1:
		return 4
	case 2:
		return 1
	case 3:
		return 1
	default:
		return 0
	}
}

func isStatus(id byte) bool {
	return id == 1
}

func decodeStatus(msg []byte) Status {
	return Status{
		moveX: int16(binary.BigEndian.Uint16(msg[:2])),
		moveY: int16(binary.BigEndian.Uint16(msg[2:])),
	}
}

func decodeAction(msg []byte) Action {
	return Action{
		value: msg[0],
	}
}
