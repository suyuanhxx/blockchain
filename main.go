package main

import (
	"github.com/satori/go.uuid"
	"strings"
	"fmt"
	"github.com/deckarep/golang-set"
	. "./core"
)

type Response struct {
	Message      string        `json:"message"`
	Index        int           `json:"index"`
	Transaction  []Transaction `json:"transaction"`
	Proof        int           `json:"proof"`
	PreviousHash string        `json:"previousHash"`
	Chain        []Block       `json:"chain"`
	NewChain     []Block       `json:"newChain"`
	TotalNodes   mapset.Set    `json:"totalNodes"`
}

var ud, _ = uuid.NewV1()
var nodeIdentifier = strings.Replace(ud.String(), "-", "", 0)

var blockchain = new(Blockchain)

func Mine() *Response {
	lastBlock := blockchain.LastBlock()
	proof := blockchain.ProofOfWork(lastBlock)

	blockchain.NewTransaction("0", nodeIdentifier, 1)

	// Forge the new Block by adding it to the chain
	previousHash := blockchain.Hash(lastBlock)
	block := blockchain.NewBlock(proof, previousHash)

	response := Response{Message: "New Block Forged",
		Index: block.Index, Transaction: block.Transactions,
		Proof: block.Proof, PreviousHash: block.PreviousHash}

	return &response
}

func NewTransaction(sender string, recipient string, amount int) *Response {
	index := blockchain.NewTransaction(sender, recipient, amount)
	response := Response{Message: "Transaction will be added to Block " + string(index), Index: index}
	return &response
}

func FullChain() *Chain {
	chain := new(Chain)
	chain.Chain = blockchain.Chain
	chain.Length = len(blockchain.Chain)
	return chain
}

func RegisterNodes(nodes []string) *Response {
	if len(nodes) <= 0 {
		return nil
	}

	for _, node := range nodes {
		blockchain.RegisterNode(node)
	}
	response := Response{Message: "New nodes have been added",
		TotalNodes: blockchain.Nodes}
	return &response
}

func Consensus() *Response {
	replaced := blockchain.ResolveConflicts()

	response := new(Response)
	if replaced {
		response.Message = "Our chain was replaced"
		response.NewChain = blockchain.Chain
	} else {
		response.Message = "Our chain is authoritative"
		response.Chain = blockchain.Chain
	}
	return response
}

func main() {
	response := Mine()
	fmt.Print(response.Message)
}
