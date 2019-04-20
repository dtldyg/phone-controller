package server

import (
	"github.com/go-vgo/robotgo"
	"fmt"
)

const (
	MOUSE_LEFT  = "left"
	MOUSE_RIGHT = "right"
	MOUSE_DOWN  = "down"
	MOUSE_UP    = "up"
)

func doAction(a Action) {
	switch a.id {
	case 2:
		fmt.Println(a.value)
		//左键按下1 左键释放2 右键按下3 右键释放4
		switch a.value {
		case 1:
			robotgo.MouseToggle(MOUSE_DOWN, MOUSE_LEFT)
		case 2:
			robotgo.MouseToggle(MOUSE_UP, MOUSE_LEFT)
		case 3:
			robotgo.MouseToggle(MOUSE_DOWN, MOUSE_RIGHT)
		case 4:
			robotgo.MouseToggle(MOUSE_UP, MOUSE_RIGHT)
		}
	case 3:
		if a.value != 0 {
			robotgo.ScrollMouse(getScroll(a))
		}
	}
}

func getScroll(a Action) (int, string) {
	if a.value < 0 {
		return -int(a.value), MOUSE_UP
	}
	return int(a.value), MOUSE_DOWN
}

func doStatus(d *Status) {
	x, y := robotgo.GetMousePos()
	robotgo.MoveMouse(x+int(d.moveX), y+int(d.moveY))
}

func doQuit() {
	robotgo.MouseToggle(MOUSE_UP, MOUSE_LEFT)
	robotgo.MouseToggle(MOUSE_UP, MOUSE_RIGHT)
}
