package server

import "sync"

func RunGo(wait *sync.WaitGroup, f func()) {
	wait.Add(1)
	go func() {
		f()
		wait.Done()
	}()
}
