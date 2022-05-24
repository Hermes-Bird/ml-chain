package wallet

import (
	"fmt"
	"github.com/Hermes-Bird/ml-chain/config"
	"github.com/Hermes-Bird/ml-chain/util"
	"log"
	"time"
)

type Transaction struct {
	Id      string
	Input   TransactionInput
	Outputs []TransactionOutput
}

type TransactionOutput struct {
	Amount  int
	Address string
}

type TransactionInput struct {
	Amount    int
	Timestamp int64
	Address   string
	Signature []byte
}

func NewTransaction() *Transaction {
	return &Transaction{
		Id:      util.Id(),
		Input:   TransactionInput{},
		Outputs: []TransactionOutput{},
	}
}

func (tx *Transaction) UpdateTransaction(senderWallet *Wallet, recipient string, amount int) {
	for i, out := range tx.Outputs {
		if out.Address == senderWallet.Address {
			if amount > out.Amount {
				log.Printf("Amount: %d exeeds balance\n", amount)
				return
			}

			out.Amount = out.Amount - amount
			tx.Outputs[i] = out

			tx.Outputs = append(tx.Outputs, TransactionOutput{
				Amount:  amount,
				Address: recipient,
			})

			SignTransaction(tx, senderWallet)

			break
		}
	}
}

func TransactionWithOutputs(senderWallet *Wallet, outputs []TransactionOutput) *Transaction {
	transaction := NewTransaction()
	transaction.Outputs = append(transaction.Outputs, outputs...)
	transaction = SignTransaction(transaction, senderWallet)
	return transaction
}

func RewardTransaction(recipient string, blockchainWallet *Wallet) *Transaction {
	return TransactionWithOutputs(blockchainWallet, []TransactionOutput{
		{Amount: config.MINING_REWARD, Address: recipient},
	})
}

func CreateTransaction(senderWallet *Wallet, recipient string, amount int) *Transaction {
	if amount > senderWallet.Balance {
		log.Printf("Amount: %d, exceeds current balance", senderWallet.Balance)
		return nil
	}

	transaction := TransactionWithOutputs(senderWallet, []TransactionOutput{
		{Amount: senderWallet.Balance - amount, Address: senderWallet.Address},
		{Amount: amount, Address: recipient},
	})

	return transaction
}

func SignTransaction(transaction *Transaction, senderWallet *Wallet) *Transaction {
	transaction.Input = TransactionInput{
		Amount:    senderWallet.Balance,
		Timestamp: time.Now().UnixMilli(),
		Address:   senderWallet.Address,
		Signature: senderWallet.Sign(util.Hash([]byte(fmt.Sprintf("%+v", transaction.Outputs)))),
	}

	return transaction
}

func VerifyTransaction(transaction *Transaction) bool {
	hash := util.Hash([]byte(fmt.Sprintf("%+v", transaction.Outputs)))
	return util.VerifySignature(transaction.Input.Address, hash, transaction.Input.Signature)
}
