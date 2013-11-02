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

// Package shamir contains the implementation of Shamir's No-Key Algorithm
package shamir

import (
	"crypto/rand"
	"math/big"
)

// GeneratePrime randomly generates a big.Int prime number with the given size
// in bits
func GeneratePrime(size int) *big.Int {
	prime, err := rand.Prime(rand.Reader, size)
	if err != nil {
		panic(err)
	}
	return prime
}

// GenerateExponents generates two exponents as big.Int, a random one and the
// inverse mod prime - 1
func GenerateExponents(prime *big.Int) (exp, expInv *big.Int) {
	primeMinusOne := big.NewInt(1).Sub(prime, big.NewInt(1))
	exp, err := rand.Int(rand.Reader, primeMinusOne)
	gcdCorrect := false
	if err == nil {
		for !gcdCorrect {
			var gcd = big.NewInt(1).GCD(nil, nil, exp, primeMinusOne)
			if gcd.Cmp(big.NewInt(1)) == 0 {
				gcdCorrect = true
			} else {
				exp = exp.Add(exp, big.NewInt(1))
				if exp.Cmp(primeMinusOne) == 0 {
					exp, err = rand.Int(rand.Reader, primeMinusOne)
				}
			}
		}
	} else {
		panic(err)
	}
	expInv = big.NewInt(1)
	expInv.ModInverse(exp, big.NewInt(1).Sub(prime, expInv))
	return
}

// Calculate encrypts or decrypts a message given as []*big.Int using the
// big.Int exponent and modulus
func Calculate(message []*big.Int, exponent *big.Int,
	modulus *big.Int) []*big.Int {
	var returnVal []*big.Int = make([]*big.Int, len(message))
	for i := 0; i < len(returnVal); i++ {
		returnVal[i] = big.NewInt(1).Exp(message[i], exponent, modulus)
	}
	return returnVal
}

// SliceMessage slices a message string into *big.Ints with a size smaller
// than the supplied prime
func SliceMessage(message string, prime *big.Int) []*big.Int {
	var returnVal []*big.Int
	var messageByteArray = []byte(message)
	if (len(messageByteArray) * 8) < prime.BitLen() {
		returnVal = make([]*big.Int, 1)
		returnVal[0] = big.NewInt(1).SetBytes(messageByteArray)
	} else {
		var size int = len(messageByteArray) * 8 / prime.BitLen()
		if len(messageByteArray)*8%prime.BitLen() > 0 {
			size++
		}
		returnVal = make([]*big.Int, size)
		var offset int = 0
		for i := 0; i < len(returnVal); i++ {
			if len(messageByteArray) > offset+(prime.BitLen()/8)+1 {
				returnVal[i] = big.NewInt(1).SetBytes(
					messageByteArray[offset : offset+(prime.BitLen()/8)])
				offset += prime.BitLen() / 8
			} else {
				returnVal[i] = big.NewInt(1).SetBytes(
					messageByteArray[offset:len(messageByteArray)])
			}
		}
	}
	return returnVal
}

// GlueMessage reassembles a message supplied as []*big.Int into one string
func GlueMessage(message []*big.Int) string {
	var preResult []byte
	for i := 0; i < len(message); i++ {
		preResult = append(preResult, message[i].Bytes()...)
	}
	return string(preResult)
}
