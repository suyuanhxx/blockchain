package core

import (
	"time"
	"math"
	"fmt"
	"strconv"
)

func (b *Blockchain) Pow(lashBlock Block) int {
	lastProof := lashBlock.Proof
	lastHash := Hash(lashBlock)
	proof := 0
	for !b.ValidProof(lastProof, proof, lastHash) {
		proof += 1
	}
	return proof
}

/**
Meet the requirements proof constituted by a string that start with n prefix '0'
 */
func (b *Blockchain) ValidProof(lastProof int, proof int, lastHash string) bool {
	guess := string(lastProof) + string(proof) + lastHash
	guessHash := Sha256(guess)
	return guessHash[:4] == "0000"
}

func Pow() int {
	timestamp, message := time.Now().String(), "This is a random message."
	nonce, guess, throttle := 0, 99999999, 100000

	payload := timestamp + message
	target := int(math.Pow(2, 6)) / throttle
	payloadHash := Sha256(payload)

	start := time.Now()
	for guess > target {
		nonce += 1
		str := Sha256(Sha256(string(nonce) + payloadHash))[0:8]
		guess, _ = strconv.Atoi(str)
		fmt.Println(guess)
	}
	end := time.Now()
	fmt.Println("consensus takes time is: ", end.Sub(start))
	return guess
}

func Pos() {

}
