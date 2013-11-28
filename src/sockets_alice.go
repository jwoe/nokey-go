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

// Package main contains all functions to demonstrate the implementation of
// Shamir's No-Key Algorithm
package main

import (
	"bufio"
	"fmt"
	"math/big"
	"net"
	"os"
	"shamir"
	"time"
)

const PRIMEBITS int = 1024

// main implements all necessary functionality to setup the conversation
// between Alice and Bob
func main() {

	msg := readMessage()
	//msg := "Hallo"
	go Client(msg)
	<-stop

}

// alice implements all the necessary functionality for Alice's part of the
// communication
func alice(msg string, channel chan []*big.Int, conn net.Conn) {
	prime := shamir.GeneratePrime(PRIMEBITS)
	primeSlice := []*big.Int{prime}
	fmt.Printf("Alice generates a prime number:\n%x\n\n",prime)

	fmt.Printf("Alice sends the prime number to Bob\n")
	channel <- primeSlice
	send(channel, conn)
	fmt.Println("Alice wants to send the following message: " + msg)
	a, aInv := shamir.GenerateExponents(prime)
	fmt.Println("Alice computes a secret Exponent and the inverse of it")
	fmt.Printf("Alice's secret exponent:\n%x\n", a)
	fmt.Printf("Alice's secret inverse:\n%x\n\n", aInv)
	fmt.Println("Alice encrypts her message!")
	var messageInt []*big.Int = shamir.SliceMessage(msg, prime)
	x := shamir.CalculateParallel(messageInt, a, prime)
	fmt.Printf("Alice now sends the encrypted message to Bob:\n%x\n\n",shamir.GlueMessage(x))
	channel <- x
	send(channel, conn)
	fmt.Println("Alice is waiting for Bob's answer...")
	receive(channel, conn)
	x = <-channel
	fmt.Println("Alice received the double-encrypted message and is now" +" decrypting her part!")
	y := shamir.CalculateParallel(x, aInv, prime)
	fmt.Printf("Alice now sends the partly decrypted message to Bob:\n%x\n\n",shamir.GlueMessage(y))
	channel <- y
	send(channel, conn)
}

func Client(msg string) {
	fmt.Printf("Client wird gestartet\n")
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:9999")
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	fmt.Printf("Verbindung erfolgreich zum Server aufgebaut\n")

	channel := make(chan []*big.Int, 1)
	go alice(msg, channel, conn)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
func receive(channel chan []*big.Int, conn net.Conn) {
	var buf = make([]byte, 1024)

	n, err := conn.Read(buf[0:])
	checkError(err)

	answer := new(big.Int)
	answer.SetBytes(buf[0:n])

	channel <- []*big.Int{answer}
}

func send(channel chan []*big.Int, conn net.Conn) {
	_, err := conn.Write([]byte(shamir.GlueMessage(<-channel)))
	checkError(err)
}

func readMessage() string {
	var message string
	fmt.Print("Please enter the message to be exchanged in encrypted form: ")
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		message = line
		if line != "" || err != nil {
			break
		}
	}
	return message
}
}