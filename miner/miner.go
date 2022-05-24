package miner

import (
	"github.com/Hermes-Bird/ml-chain/blockchain"
	"github.com/Hermes-Bird/ml-chain/p2p"
	"github.com/Hermes-Bird/ml-chain/wallet"
)

type Miner struct {
	bc        *blockchain.Blockchain
	tp        *wallet.TransactionPool
	wall      *wallet.Wallet
	p2pServer *p2p.P2PServer
}

func NewMiner(bc *blockchain.Blockchain, tp *wallet.TransactionPool, wall *wallet.Wallet, p2pServer *p2p.P2PServer) *Miner {
	return &Miner{
		bc:        bc,
		tp:        tp,
		wall:      wall,
		p2pServer: p2pServer,
	}
}

func (m Miner) Mine() *blockchain.Block {
	validTransactions := m.tp.ValidTransactions()
	validTransactions = append(validTransactions, *wallet.RewardTransaction(m.wall.Address, wallet.BlockchainWallet()))
	block := m.bc.AddBlock(validTransactions)
	m.p2pServer.SyncChains()
	m.tp.Clear()
	m.p2pServer.BroadcastClear()
	return &block
}
