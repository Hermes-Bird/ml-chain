package blockchain

import (
	"fmt"
	"github.com/Hermes-Bird/ml-chain/wallet"
	"log"
)

type Blockchain struct {
	Chain []Block
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		Chain: []Block{GenesisBlock()},
	}
}

func (bc *Blockchain) AddBlock(data []wallet.Transaction) Block {
	newBlock := MineBlock(bc.Chain[len(bc.Chain)-1], data)
	bc.Chain = append(bc.Chain, newBlock)

	return newBlock
}

func (bc *Blockchain) IsValidChain(chain []Block) bool {
	if fmt.Sprint(chain[0]) != fmt.Sprint(GenesisBlock()) {
		return false
	}

	for i := 1; i < len(chain); i++ {
		prevBlock := chain[i-1]
		currentBlock := chain[i]
		if currentBlock.LastHash != prevBlock.Hash || BlockHash(currentBlock) != currentBlock.Hash {
			return false
		}
	}

	return true
}

func (bc *Blockchain) ReplaceChain(newChain []Block) {
	if len(bc.Chain) >= len(newChain) {
		log.Println("New chain isn't longer than the current chain")
		return
	}
	if !bc.IsValidChain(newChain) {
		log.Println("Invalid chain received")
		return
	}

	log.Println("Replace current chain with a new one")
	bc.Chain = newChain
}

func (bc *Blockchain) GetChainTransactions() []wallet.Transaction {
	var transactions []wallet.Transaction

	for _, block := range bc.Chain {
		for _, tx := range block.Data {
			transactions = append(transactions, tx)
		}
	}

	return transactions
}
