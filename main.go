package main

import (
	"github.com/satori/go.uuid"
	"strings"
	. "./core"
	"time"
	"net/http"
	"log"
	"encoding/json"
)

type Response struct {
	Message      string        `json:"message"`
	Index        int           `json:"index"`
	Transaction  []Transaction `json:"transaction"`
	Proof        int           `json:"proof"`
	PreviousHash string        `json:"previousHash"`
	Chain        []Block       `json:"chain,omitempty"`
	NewChain     []Block       `json:"newChain,omitempty"`
	TotalNodes   []interface{} `json:"totalNodes,omitempty"`
}

var ud, _ = uuid.NewV1()
var nodeIdentifier = strings.Replace(ud.String(), "-", "", -1)

var blockchain *Blockchain

func mine(w http.ResponseWriter, r *http.Request) {
	lastBlock := blockchain.LastBlock()
	proof := blockchain.Pow(lastBlock)

	blockchain.NewTransaction("0", nodeIdentifier, 1)

	// Forge the new Block by adding it to the chain
	previousHash := Hash(lastBlock)
	block := blockchain.NewBlock(proof, previousHash)

	response := Response{Message: "New Block Forged",
		Index: block.Index, Transaction: block.Transactions,
		Proof: block.Proof, PreviousHash: block.PreviousHash}

	resp, _ := json.Marshal(response)

	w.Write(resp)

}

func NewTransaction(sender string, recipient string, amount int) *Response {
	index := blockchain.NewTransaction(sender, recipient, amount)
	response := Response{Message: "Transaction will be added to Block " + string(index), Index: index}
	return &response
}

func fullChain(w http.ResponseWriter, r *http.Request) {
	chain := new(Chain)
	chain.Chain = blockchain.Chain
	chain.Length = len(blockchain.Chain)

	resp, _ := json.Marshal(chain)
	w.Write(resp)
}

func RegisterNodes(nodes []string) *Response {
	if len(nodes) <= 0 {
		return nil
	}

	for _, node := range nodes {
		blockchain.RegisterNode(node)
	}

	response := Response{Message: "New nodes have been added",
		TotalNodes: blockchain.Nodes.ToSlice()}
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

type HandlersFunc func(http.ResponseWriter, *http.Request)

var handlersMap = make(map[string]HandlersFunc)

type Server struct {
}

func (*Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := handlersMap[r.URL.String()]; ok {
		h(w, r)
	}
}

func main() {
	blockchain = blockchain.New()

	server := &http.Server{
		Addr:           ":8000",
		Handler:        &Server{},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	handlersMap["/mine"] = mine
	handlersMap["/chain"] = fullChain
	log.Fatal(server.ListenAndServe())

}
