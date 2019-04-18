package server

import (
	"github.com/go-vgo/robotgo"
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

func doStatus(d *Data) {
	//int16[-32700-32700] -> float64[0-327]
	speedX := float64(d.status.speedX) / 100
	speedY := float64(d.status.speedY) / 100

	x, y := robotgo.GetMousePos()
	if absFloat64(float64(x)-d.xF) >= 1 {
		d.xF = float64(x)
	}
	if absFloat64(float64(y)-d.yF) >= 1 {
		d.yF = float64(y)
	}
	d.xF += speedX
	d.yF -= speedY //y向下为正
	robotgo.MoveMouse(int(d.xF), int(d.yF))
}

func doQuit() {
	robotgo.MouseToggle(MOUSE_UP, MOUSE_LEFT)
	robotgo.MouseToggle(MOUSE_UP, MOUSE_RIGHT)
}

func absFloat64(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}
