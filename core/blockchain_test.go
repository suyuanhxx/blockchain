package core

import (
	"testing"
	"fmt"
)

func Test_RegisterNode(t *testing.T) {
	blockchain := new(Blockchain)
	blockchain.RegisterNode("1111")
}

func Test_Pow(t *testing.T) {
	n := Pow()
	fmt.Print(n)
}
