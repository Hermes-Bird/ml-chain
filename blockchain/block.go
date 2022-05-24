package blockchain

import (
	"encoding/hex"
	"fmt"
	"github.com/Hermes-Bird/ml-chain/config"
	"github.com/Hermes-Bird/ml-chain/util"
	"github.com/Hermes-Bird/ml-chain/wallet"
	"strings"
	"time"
)

type Block struct {
	Difficulty int
	Timestamp  int64
	LastHash   string
	Hash       string
	Nonce      int
	Data       []wallet.Transaction
}

func GenesisBlock() Block {
	return Block{
		Difficulty: config.DIFFICULTY,
		LastHash:   "LastGenesisHash",
		Hash:       "GenesisHash",
		Timestamp:  0,
		Nonce:      0,
	}
}
func Hash(timestamp int64, lastHash string, data any, nonce int, difficulty int) string {
	hashBytes := util.Hash([]byte(fmt.Sprintf("{%v}{%v}{%v}{%v}${%v}", timestamp, lastHash, data, nonce, difficulty)))
	return hex.EncodeToString(hashBytes)
}

func BlockHash(block Block) string {
	return Hash(block.Timestamp, block.LastHash, block.Data, block.Nonce, block.Difficulty)
}

func MineBlock(lastBlock Block, data []wallet.Transaction) Block {
	var timestamp int64
	lastHash := lastBlock.Hash
	difficulty := lastBlock.Difficulty
	nonce := 0

	hash := ""
	for {
		timestamp = time.Now().Unix() * 1000
		difficulty = adjustDifficulty(lastBlock, timestamp)
		hash = Hash(timestamp, lastHash, data, nonce, difficulty)
		if strings.HasPrefix(hash, strings.Repeat("0", difficulty)) {
			break
		} else {
			nonce++
		}
	}
	return NewBlock(timestamp, lastHash, hash, data, nonce, difficulty)
}

func adjustDifficulty(lastBlock Block, timestamp int64) int {
	if lastBlock.Timestamp+config.MINE_RATE > timestamp {
		return lastBlock.Difficulty + 1
	}

	return lastBlock.Difficulty - 1
}

func NewBlock(timestamp int64, lastHash string, hash string, data []wallet.Transaction, nonce int, difficulty int) Block {
	return Block{
		Difficulty: difficulty,
		Timestamp:  timestamp,
		LastHash:   lastHash,
		Hash:       hash,
		Data:       data,
		Nonce:      nonce,
	}
}
