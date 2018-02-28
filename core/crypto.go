package core

import (
	"encoding/hex"
	"crypto/sha256"
)

var hash256 = sha256.New()

func ValidProof(lastProof int, proof int, lastHash string) bool {
	guess := string(lastProof) + string(proof) + lastHash
	guessHash := Sha256(guess)
	return guessHash[:4] == "0000"
}

func Sha256(str string) string {
	hash256.Write([]byte(str))
	return hex.EncodeToString(hash256.Sum(nil))
}
