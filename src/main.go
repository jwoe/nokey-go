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
	"bufio"
	"fmt"
	"math/big"
	"os"
	"shamir"
)

const PRIMEBITS int = 1024

// main implements all necessary functionality to setup the conversation
// between Alice and Bob
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
	prime := shamir.GeneratePrime(PRIMEBITS)
	fmt.Printf("Both Alice and Bob agree on a prime number:\n%x\n\n",
		prime)
	channel := make(chan []*big.Int)
	stop := make(chan int)

	go alice(message, prime, channel)
	go bob(prime, channel, stop)
	<-stop
}

// alice implements all the necessary functionality for Alice's part of the
// communication
func alice(msg string, prime *big.Int, channel chan []*big.Int) {
	fmt.Println("Alice wants to send the following message: " + msg)
	a, aInv := shamir.GenerateExponents(prime)
	fmt.Println("Alice computes a secret Exponent and the inverse of it")
	fmt.Printf("Alice's secret exponent:\n%x\n", a)
	fmt.Printf("Alice's secret inverse:\n%x\n\n", aInv)
	fmt.Println("Alice encrypts her message!")
	var messageInt []*big.Int = shamir.SliceMessage(msg, prime)
	x := shamir.Calculate(messageInt, a, prime)
	fmt.Printf("Alice now sends the encrypted message to Bob:\n%x\n\n",
		shamir.GlueMessage(x))
	channel <- x
	fmt.Println("Alice is waiting for Bob's answer...")
	x = <-channel
	fmt.Println("Alice received the double-encrypted message and is now" +
		" decrypting her part!")
	y := shamir.Calculate(x, aInv, prime)
	fmt.Printf("Alice now sends the partly decrypted message to Bob:\n%x\n\n",
		shamir.GlueMessage(y))
	channel <- y
}

// bob implements all the necessary functionality for Bob's part of the
// communication
func bob(prime *big.Int, channel chan []*big.Int, stop chan int) {
	fmt.Println("Bob is waiting for the encrypted message from Alice...")
	x := <-channel
	b, bInv := shamir.GenerateExponents(prime)
	fmt.Println("Bob computes a secret Exponent and the inverse of it")
	fmt.Printf("Bob's secret exponent:\n%x\n", b)
	fmt.Printf("Bob's secret inverse:\n%x\n\n", bInv)
	fmt.Println("Bob received the encrypted message from Alice and is now" +
		" encrypting it too!")
	y := shamir.Calculate(x, b, prime)
	fmt.Printf("Bob now sends the double-encrypted message back to "+
		"Alice:\n%x\n\n", shamir.GlueMessage(y))
	channel <- y
	fmt.Println("Bob is waiting for Alice's answer...")
	x = <-channel
	fmt.Println("Bob received the second message from Alice and is now " +
		"decrypting it!")
	y = shamir.Calculate(x, bInv, prime)
	fmt.Println("Bob decrypted the following message from Alice: " +
		shamir.GlueMessage(y))
	stop <- 1
}
