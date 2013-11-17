/* Copyright 2013 Georg Hartmann, Manuel Schweizer

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License. */

package main

import (
	//	"bufio"
	"fmt"
	"math"
	"math/big"
	//	"os"
	"bytes"
	"runtime"
	"shamir"
	"time"
)

//const PRIMEBITS int = 1024

// main implements all necessary functionality to setup the conversation
// between Alice and Bob
func runtimeTest() {
	fmt.Println("CPU 's # ", runtime.NumCPU())

	prime := shamir.GeneratePrime(PRIMEBITS)

	var message string

	//fmt.Print(prime, message)

	a, aInv := shamir.GenerateExponents(prime)

	var time0 time.Time
	var duration time.Duration

	//var words int
	for nCpu := 1; nCpu <= runtime.NumCPU(); nCpu++ {
		runtime.GOMAXPROCS(nCpu)
		fmt.Println("Using CPU # ", nCpu)
		//fmt.Println("\t\t Dauer Normal \t\tDauer Parallel")
		for i := 1; i < 5; i++ {
			word := math.Pow10(i)
			//	words = 10 ^ i
			var buffer bytes.Buffer
			fmt.Print("Chars # ", word*50, "\t\t")
			var j float64 = 0
			for j = 0; j < word; j++ {
				buffer.WriteString("abcdefghijklmnopqrstuvwxzy 1234567890?!, .-#+* ()[]")
			}
			message = buffer.String()

			var messageInt []*big.Int = shamir.SliceMessage(message, prime)

			// Normal
			time0 = time.Now()
			shamir.Calculate(messageInt, a, prime)
			duration = time.Since(time0)

			fmt.Print(float64(duration.Nanoseconds()) / 1000 / 1000)

			// Parrallel
			time0 = time.Now()
			shamir.CalculateParallel(messageInt, a, prime)
			duration = time.Since(time0)
			fmt.Print("\t\t", float64(duration.Nanoseconds())/1000/1000)
			fmt.Println("")

		}

	}

	var messageInt []*big.Int = shamir.SliceMessage("foo", prime)

	x := shamir.CalculateParallel(messageInt, a, prime)
	shamir.GlueMessage(x)

	shamir.Calculate(x, aInv, prime)

}
