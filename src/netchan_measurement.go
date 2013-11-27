// netchan_measurement.go
package main

import (
	"code.google.com/p/go.exp/old/netchan"
	"fmt"
	"time"
)

func main() {
	msg := "Hallo"
	stop := make(chan int)

	var time0 time.Time
	var duration time.Duration
	time0 = time.Now()

	go Server(stop)
	go Client(msg)
	<-stop
	duration = time.Since(time0)

	fmt.Print(float64(duration.Nanoseconds()) / 1000 / 1000)
}

func Server(stop chan int) {
	exp := netchan.NewExporter()

	errNE := exp.ListenAndServe("tcp", ":9999")
	handleError(errNE)
	channelBobReceive := make(chan string)
	errRE := exp.Export("AliceToBob", channelBobReceive, netchan.Recv)
	handleError(errRE)

	for i := 0; i < 100; i++ {
		<-channelBobReceive
		//fmt.Printf("%i %s\n", i, <-channelBobReceive)
	}
	stop <- 1
}

func Client(msg string) {
	imp, err := netchan.Import("tcp", "127.0.0.1:9999")
	handleError(err)

	channelAliceSend := make(chan string)
	imp.Import("AliceToBob", channelAliceSend, netchan.Send, 1)
	for i := 1; i < 101; i++ {
		channelAliceSend <- msg
	}
}
func handleError(err error) {
	if err != nil {
		fmt.Println("Error: %v", err)
	}
}
