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

// netchan_Alice
package main

import (
	"bufio"
	"code.google.com/p/go.exp/old/netchan"
	"fmt"
	"math/big"
	"os"
	"shamir"
)

const PRIMEBITS int = 1024

func main() {
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

	imp, err := netchan.Import("tcp", "127.0.0.1:9999")
	handleError(err)

	channelAliceSend := make(chan []*big.Int)
	channelAliceReceive := make(chan []*big.Int)

	imp.Import("BobToAlice", channelAliceReceive, netchan.Recv, 1)
	fmt.Println("Empfangskanal importiert")
	imp.Import("AliceToBob", channelAliceSend, netchan.Send, 1)
	fmt.Println("Sendekanal importiert")

	stop := make(chan int)

	go alice(message, channelAliceReceive, channelAliceSend, stop)
	<-stop
}

// alice implements all the necessary functionality for Alice's part of the
// communication
func alice(msg string, channelReceive chan []*big.Int, channelSend chan []*big.Int, stop chan int) {
	prime := shamir.GeneratePrime(PRIMEBITS)
	primeSlice := []*big.Int{prime}
	fmt.Printf("Alice generates a prime number:\n%x\n\n",
		prime)

	fmt.Printf("Alice sends the prime number to Bob\n")
	channelSend <- primeSlice
	fmt.Println("Alice wants to send the following message: " + msg)
	a, aInv := shamir.GenerateExponents(prime)
	fmt.Println("Alice computes a secret Exponent and the inverse of it")
	fmt.Printf("Alice's secret exponent:\n%x\n", a)
	fmt.Printf("Alice's secret inverse:\n%x\n\n", aInv)
	fmt.Println("Alice encrypts her message!")
	var messageInt []*big.Int = shamir.SliceMessage(msg, prime)
	//x := shamir.Calculate(messageInt, a, prime)
	x := shamir.CalculateParallel(messageInt, a, prime)
	fmt.Printf("Alice now sends the encrypted message to Bob:\n%x\n\n",
		shamir.GlueMessage(x))
	channelSend <- x
	fmt.Println("Alice is waiting for Bob's answer...")
	x = <-channelReceive
	fmt.Println("Alice received the double-encrypted message and is now" +
		" decrypting her part!")
	y := shamir.CalculateParallel(x, aInv, prime)
	fmt.Printf("Alice now sends the partly decrypted message to Bob:\n%x\n\n",
		shamir.GlueMessage(y))
	channelSend <- y
	stop <- 1
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Error: %v", err)
	}
}
