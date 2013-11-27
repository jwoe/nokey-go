/* Copyright 2013 Michael Galetzka, Jonas Woerlein

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
	//"bufio"
	"code.google.com/p/go.exp/old/netchan"
	"fmt"
	"math/big"
	//"os"
	"shamir"
	"time"
)

const PRIMEBITS int = 1024

// main implements all necessary functionality to setup the conversation
// between Alice and Bob
func main() {
	/*var message string
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
	}*/
	message := "Hallo"

	var time0 time.Time
	var duration time.Duration
	time0 = time.Now()

	exp := netchan.NewExporter()

	errNE := exp.ListenAndServe("tcp", ":9999")
	handleError(errNE)

	imp, err := netchan.Import("tcp", "127.0.0.1:9999")
	handleError(err)

	channelAliceSend := make(chan []*big.Int)
	channelAliceReceive := make(chan []*big.Int)
	channelBobSend := make(chan []*big.Int)
	channelBobReceive := make(chan []*big.Int)

	errSE := exp.Export("BobToAlice", channelBobSend, netchan.Send)
	handleError(errSE)

	errRE := exp.Export("AliceToBob", channelBobReceive, netchan.Recv)
	handleError(errRE)

	imp.Import("BobToAlice", channelAliceReceive, netchan.Recv, 1)
	//fmt.Println("Empfangskanal importiert")
	imp.Import("AliceToBob", channelAliceSend, netchan.Send, 1)
	//fmt.Println("Sendekanal importiert")

	stop := make(chan int)
	for i := 1; i < 100; i++ {
		go alice(message, channelAliceReceive, channelAliceSend)
		go bob(channelBobReceive, channelBobSend, stop)
		<-stop
	}
	duration = time.Since(time0)

	fmt.Print(float64(duration.Nanoseconds()) / 1000 / 1000)
}

// alice implements all the necessary functionality for Alice's part of the
// communication
func alice(msg string, channelReceive chan []*big.Int, channelSend chan []*big.Int) {
	prime := shamir.GeneratePrime(PRIMEBITS)
	primeSlice := []*big.Int{prime}
	//fmt.Printf("Alice generates a prime number:\n%x\n\n",prime)

	//fmt.Printf("Alice sends the prime number to Bob\n")
	channelSend <- primeSlice
	//fmt.Println("Alice wants to send the following message: " + msg)
	a, aInv := shamir.GenerateExponents(prime)
	//fmt.Println("Alice computes a secret Exponent and the inverse of it")
	//fmt.Printf("Alice's secret exponent:\n%x\n", a)
	//fmt.Printf("Alice's secret inverse:\n%x\n\n", aInv)
	//fmt.Println("Alice encrypts her message!")
	var messageInt []*big.Int = shamir.SliceMessage(msg, prime)
	x := shamir.CalculateParallel(messageInt, a, prime)
	//fmt.Printf("Alice now sends the encrypted message to Bob:\n%x\n\n",shamir.GlueMessage(x))
	channelSend <- x
	//fmt.Println("Alice is waiting for Bob's answer...")
	x = <-channelReceive
	//fmt.Println("Alice received the double-encrypted message and is now" + " decrypting her part!")
	y := shamir.CalculateParallel(x, aInv, prime)
	//fmt.Printf("Alice now sends the partly decrypted message to Bob:\n%x\n\n", shamir.GlueMessage(y))
	channelSend <- y
}

// bob implements all the necessary functionality for Bob's part of the
// communication
func bob(channelReceive chan []*big.Int, channelSend chan []*big.Int, stop chan int) {
	//fmt.Printf("Bob is waiting for a prime number from Alice...")
	primeSlice := <-channelReceive

	prime := primeSlice[0]
	if !(*prime).ProbablyPrime(4) {
		fmt.Printf("Alice prime number is probably not prime")
	}
	//fmt.Println("Bob is waiting for the encrypted message from Alice...")
	x := <-channelReceive
	b, bInv := shamir.GenerateExponents(prime)
	//fmt.Println("Bob computes a secret Exponent and the inverse of it")
	//fmt.Printf("Bob's secret exponent:\n%x\n", b)
	//fmt.Printf("Bob's secret inverse:\n%x\n\n", bInv)
	//fmt.Println("Bob received the encrypted message from Alice and is now" +" encrypting it too!")
	y := shamir.CalculateParallel(x, b, prime)
	//fmt.Printf("Bob now sends the double-encrypted message back to "+"Alice:\n%x\n\n", shamir.GlueMessage(y))
	channelSend <- y
	//fmt.Println("Bob is waiting for Alice's answer...")
	x = <-channelReceive
	//fmt.Println("Bob received the second message from Alice and is now " +"decrypting it!")
	y = shamir.CalculateParallel(x, bInv, prime)
	//fmt.Println("Bob decrypted the following message from Alice: " +shamir.GlueMessage(y))
	stop <- 1
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Error: %v", err)
	}
}
