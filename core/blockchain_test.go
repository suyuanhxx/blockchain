package core

import "testing"
import 	. "./"

func Test_RegisterNode(t *testing.T) {
	blockchain := new(Blockchain)
	blockchain.RegisterNode("1111")
}
