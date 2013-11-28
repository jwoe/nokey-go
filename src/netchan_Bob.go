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

// netchan_Bob
package main

import (
	"code.google.com/p/go.exp/old/netchan"
	"fmt"
	"math/big"
	"shamir"
)

func main() {
	exp := netchan.NewExporter()

	errNE := exp.ListenAndServe("tcp", ":9999")
	handleError(errNE)

	channelBobSend := make(chan []*big.Int)
	channelBobReceive := make(chan []*big.Int)

	errSE := exp.Export("BobToAlice", channelBobSend, netchan.Send)
	handleError(errSE)

	errRE := exp.Export("AliceToBob", channelBobReceive, netchan.Recv)
	handleError(errRE)

	stop := make(chan int)

	go bob(channelBobReceive, channelBobSend, stop)
	<-stop
}

// bob implements all the necessary functionality for Bob's part of the
// communication
func bob(channelReceive chan []*big.Int, channelSend chan []*big.Int, stop chan int) {
	fmt.Printf("Bob is waiting for a prime number from Alice...")
	primeSlice := <-channelReceive

	prime := primeSlice[0]
	if !(*prime).ProbablyPrime(4) {
		fmt.Printf("Alice prime number is probably not prime")
	}
	fmt.Println("Bob is waiting for the encrypted message from Alice...")
	x := <-channelReceive
	b, bInv := shamir.GenerateExponents(prime)
	fmt.Println("Bob computes a secret Exponent and the inverse of it")
	fmt.Printf("Bob's secret exponent:\n%x\n", b)
	fmt.Printf("Bob's secret inverse:\n%x\n\n", bInv)
	fmt.Println("Bob received the encrypted message from Alice and is now" +
		" encrypting it too!")
	y := shamir.CalculateParallel(x, b, prime)
	fmt.Printf("Bob now sends the double-encrypted message back to "+
		"Alice:\n%x\n\n", shamir.GlueMessage(y))
	channelSend <- y
	fmt.Println("Bob is waiting for Alice's answer...")
	x = <-channelReceive
	fmt.Println("Bob received the second message from Alice and is now " +
		"decrypting it!")
	y = shamir.CalculateParallel(x, bInv, prime)
	fmt.Println("Bob decrypted the following message from Alice: " +
		shamir.GlueMessage(y))
	stop <- 1
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Error: %v", err)
	}
}
