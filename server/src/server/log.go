package server

import (
	"time"
	"fmt"
	"os"
)

func LogInfo(format string, a ...interface{}) {
	fmt.Printf("[%s] %v\n", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, a...))
}

func LogError(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "[%s] %v\n", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, a...))
}
