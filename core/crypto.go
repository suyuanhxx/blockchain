package core

import (
	"encoding/hex"
	"crypto/sha256"
	"encoding/json"
)

var hash256 = sha256.New()

func Hash(block Block) string {
	blockString, _ := json.Marshal(block)
	return Sha256(string(blockString))
}

func Sha256(str string) string {
	hash256.Write([]byte(str))
	return hex.EncodeToString(hash256.Sum(nil))
}
