package main

import (
	"fmt"
	"time"
	"github.com/go-vgo/robotgo"
)

func main() {
	fmt.Println("begin")
	time.Sleep(time.Second * 3)
	//robotgo.MouseToggle("down", "left")
	//time.Sleep(time.Second * 5)
	//robotgo.MouseToggle("up", "left")

	fmt.Println(robotgo.GetScreenSize())

	//0-1
}
