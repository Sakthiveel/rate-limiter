package main

import (
	"fmt"
	"time"
)

func main() {
	userInputSecond := 2
	seconds := userInputSecond * int(time.Second)
	beforeTime := time.Now()
	fmt.Println("Before time", beforeTime)
	time.Sleep(time.Duration(seconds))
	// blockedUntil := time.Now().Add(time.Duration(seconds))
	blockedUntil := time.Now()
	timeDiff := blockedUntil.Sub(beforeTime).Seconds()
	fmt.Println("Before time", beforeTime)
	fmt.Println("After time", blockedUntil)
	fmt.Println("Diff:", timeDiff)
}
