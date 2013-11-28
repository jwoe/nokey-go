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

// socket_measurement.go
package main

import (
	"fmt"
	"net"
	"os"
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

func Client(msg string) {
	//fmt.Printf("Client wird gestartet\n")
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "192.168.0.112:9999")
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	//fmt.Printf("Verbindung erfolgreich zum Server aufgebaut\n")
	for i := 0; i < 100; i++ {
		_, err := conn.Write([]byte(msg))
		checkError(err)
	}

}

func Server(stop chan int) {
	//fmt.Printf("Server wird gestartet\n")
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "192.168.0.112:9999")
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	//fmt.Printf("Server wartet auf Verbindungsanfragen\n")

	conn, err := listener.Accept()
	if err != nil {

	}

	var buf = make([]byte, 5)

	for i := 0; i < 100; i++ {

		_, err := conn.Read(buf[0:])
		checkError(err)
		//fmt.Printf("%i %s\n", i, buf[0:n])
	}

	listener.Close()
	stop <- 1
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
