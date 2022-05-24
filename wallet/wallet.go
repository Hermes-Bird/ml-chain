package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"github.com/Hermes-Bird/ml-chain/config"
	"github.com/Hermes-Bird/ml-chain/util"
	"log"
)

type Wallet struct {
	Balance    int
	PrivateKey *ecdsa.PrivateKey
	Address    string
}

func NewWallet() *Wallet {
	privateKey, err := util.GenKeyPair()
	if err != nil {
		log.Fatalln(err)
	}
	return &Wallet{
		Balance:    config.INITIAL_BALANCE,
		PrivateKey: privateKey,
		Address:    util.GenWalletAddress(privateKey),
	}
}

func (w *Wallet) Sign(dataHash []byte) []byte {
	sign, err := ecdsa.SignASN1(rand.Reader, w.PrivateKey, dataHash)
	if err != nil {
		return []byte{}
	}
	return sign
}

func (w *Wallet) CreateTransaction(recipient string, amount int, txs []Transaction, txPool *TransactionPool) *Transaction {
	w.Balance = w.CalculateBalance(txs)

	if amount > w.Balance {
		log.Printf("Amount: %d exeeds current balance\n")
	}

	tx := txPool.ExistingTx(w.Address)

	if tx != nil {
		tx.UpdateTransaction(w, recipient, amount)
	} else {
		tx = CreateTransaction(w, recipient, amount)
		txPool.AddOrUpdateTx(tx)
	}

	return tx
}

func (w *Wallet) CalculateBalance(transactions []Transaction) int {
	balance := w.Balance

	var walletInputTxs []Transaction
	for _, tx := range transactions {
		if tx.Input.Address == w.Address {
			walletInputTxs = append(walletInputTxs, tx)
		}
	}

	var startTime int64 = 0

	if len(walletInputTxs) > 0 {
		var minTimeTransaction Transaction

		for _, tx := range walletInputTxs {
			if tx.Input.Timestamp > minTimeTransaction.Input.Timestamp {
				minTimeTransaction = tx
			}
		}

		startTime = minTimeTransaction.Input.Timestamp

		for _, out := range minTimeTransaction.Outputs {
			if out.Address == w.Address {
				balance = out.Amount
				break
			}
		}
	}

	for _, tx := range transactions {
		if tx.Input.Timestamp > startTime {
			for _, out := range tx.Outputs {
				balance += out.Amount
			}
		}
	}

	return balance
}

func BlockchainWallet() *Wallet {
	priv, _ := util.GenKeyPair()
	return &Wallet{
		PrivateKey: priv,
		Address:    "blockchain-wallet",
	}
}
