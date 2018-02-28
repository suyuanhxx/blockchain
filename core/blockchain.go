package core

import (
	"encoding/json"
	"time"
	"net/url"
	"fmt"
	"net/http"
	"github.com/deckarep/golang-set"
)

type Block struct {
	Index        int           `json:"index"`
	Timestamp    time.Time     `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	Proof        int           `json:"proof"`
	PreviousHash string        `json:"previous_hash"`
}

type Transaction struct {
	Amount    int    `json:"amount"`
	Recipient string `json:"recipient"`
	Sender    string `json:"sender"`
}

type Chain struct {
	Chain  []Block `json:"chain"`
	Length int     `json:"length"`
}

type Blockchain struct {
	CurrentTransactions []Transaction
	Chain               []Block
	Nodes               mapset.Set
}

func (t *Blockchain) RegisterNode(address string) {
	parsedUrl, _ := url.Parse(address)
	if t.Nodes == nil {
		t.Nodes = mapset.NewSet()
	}
	t.Nodes.Add(parsedUrl.Host)
}

func (t *Blockchain) ValidChain(chain []Block) bool {
	lastBlock := chain[0]
	currentIndex := 1

	for ; currentIndex < len(chain); {
		block := chain[currentIndex]
		fmt.Println(lastBlock, block)
		if block.PreviousHash != t.Hash(lastBlock) {
			return false
		}

		if !ValidProof(lastBlock.Proof, block.Proof, lastBlock.PreviousHash) {
			return false
		}
		lastBlock = block
		currentIndex += 1
	}
	return true
}

func (t *Blockchain) ResolveConflicts() bool {
	var newChain []Block

	maxLength := len(t.Chain)

	for node := range t.Nodes.Iter() {
		resp, error := http.Get("http://" + node.(string) + "/chain")
		if error != nil {
			continue
		}
		defer resp.Body.Close()
		var chain Chain
		json.NewDecoder(resp.Body).Decode(&chain)

		if chain.Length > maxLength && t.ValidChain(chain.Chain) {
			maxLength = chain.Length
			newChain = chain.Chain
		}
	}

	if newChain != nil {
		t.Chain = newChain
		return true
	}
	return false

}

func (t *Blockchain) NewBlock(proof int, previousHash string) Block {
	block := new(Block)
	block.Index = len(t.Chain) + 1
	block.Timestamp = time.Now()
	block.Proof = proof
	if previousHash != "" {
		block.PreviousHash = previousHash
	} else {
		block.PreviousHash = t.Hash(t.Chain[len(t.Chain)-1])
	}
	//t.CurrentTransactions = []
	t.Chain = append(t.Chain, *block)
	return *block
}

func (t *Blockchain) NewTransaction(sender string, recipient string, amount int) int {
	transaction := Transaction{Sender: sender, Recipient: recipient, Amount: amount}
	t.CurrentTransactions = append(t.CurrentTransactions, transaction)
	return t.LastBlock().Index + 1
}

func (t *Blockchain) LastBlock() Block {
	return t.Chain[len(t.Chain)-1]
}

func (t *Blockchain) ProofOfWork(lashBlock Block) int {
	lastProof := lashBlock.Proof
	lastHash := t.Hash(lashBlock)
	proof := 0
	for !ValidProof(lastProof, proof, lastHash) {
		proof += 1
	}
	return 0
}

func (t *Blockchain) Hash(block Block) string {
	blockString, _ := json.Marshal(block)
	return Sha256(string(blockString))
}

//func ValidProof(lastProof int, proof int, lastHash string) bool {
//	guess := string(lastProof) + string(proof) + lastHash
//	guessHash := Sha256(guess)
//	return guessHash[:4] == "0000"
//}
//
//func Sha256(str string) string {
//	hash256.Write([]byte(str))
//	return hex.EncodeToString(hash256.Sum(nil))
//}
