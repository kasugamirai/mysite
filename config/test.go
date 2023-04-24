package config

import (
	"fmt"
	"time"
)

func sendData(ch chan int) {
	for i := 0; i < 5; i++ {
		ch <- i
		time.Sleep(time.Millisecond * 100)
	}
	close(ch)
}

func main() {
	dataChannel := make(chan int)
	go sendData(dataChannel)

	for v := range dataChannel {
		fmt.Println("Received:", v)
	}
}
