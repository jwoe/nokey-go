/* Copyright 2013 Michael Galetzka, Jonas Woerlein, Christoph Tonnier

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License. */

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
