package app

import (
	"github.com/Hermes-Bird/ml-chain/api"
	"github.com/Hermes-Bird/ml-chain/blockchain"
	"github.com/Hermes-Bird/ml-chain/miner"
	"github.com/Hermes-Bird/ml-chain/p2p"
	"github.com/Hermes-Bird/ml-chain/wallet"
	"os"
	"strings"
)

func Start() {
	wsPort := os.Getenv("WS_PORT")
	apiPort := os.Getenv("HTTP_PORT")
	envPeers := os.Getenv("PEERS")

	peers := strings.Split(envPeers, ",")

	chain := blockchain.NewBlockchain()
	wall := wallet.NewWallet()
	txPool := wallet.NewTransactionPool()

	p2pServer := p2p.NewP2PServer(chain, txPool)
	mr := miner.NewMiner(chain, txPool, wall, p2pServer)
	go p2pServer.Start(wsPort)
	p2pServer.ConnectPeers(peers)
	api.Start(apiPort, chain, p2pServer, wall, txPool, mr)
}
