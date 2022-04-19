package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Block struct {
	Timestamp int64
	LastHash  string
	Hash      string
	Data      any
}

func GenesisBlock() Block {
	return Block{
		Timestamp: 0,
		LastHash:  "LastGenesisHash",
		Hash:      "GenesisHash",
	}
}

func Hash(timestamp int64, lastHash string, data any) string {
	hashBytes := sha256.Sum256([]byte(fmt.Sprintf("{%v}{%v}{%v}", timestamp, lastHash, data)))
	return hex.EncodeToString(hashBytes[:])
}

func BlockHash(block Block) string {
	return Hash(block.Timestamp, block.LastHash, block.Data)
}

func MineBlock(lastBlock Block, data any) Block {
	timestamp := time.Now().Unix()
	lastHash := lastBlock.Hash
	hash := Hash(timestamp, lastHash, data)

	return NewBlock(timestamp, lastHash, hash, data)
}

func NewBlock(timestamp int64, lastHash string, hash string, data any) Block {
	return Block{
		Timestamp: timestamp,
		LastHash:  lastHash,
		Hash:      hash,
		Data:      data,
	}
}
