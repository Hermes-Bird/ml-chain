package wallet

import (
	"log"
)

type TransactionPool struct {
	Transactions []Transaction
}

func NewTransactionPool() *TransactionPool {
	return &TransactionPool{Transactions: []Transaction{}}
}

func (txPool *TransactionPool) AddOrUpdateTx(transaction *Transaction) {
	for i, tx := range txPool.Transactions {
		if tx.Id == transaction.Id {
			txPool.Transactions[i] = *transaction
			return
		}
	}

	txPool.Transactions = append(txPool.Transactions, *transaction)
}

func (txPool TransactionPool) ExistingTx(address string) *Transaction {
	for i := range txPool.Transactions {
		if txPool.Transactions[i].Input.Address == address {
			return &txPool.Transactions[i]
		}
	}
	return nil
}

func (txPool *TransactionPool) ValidTransactions() []Transaction {
	var validTxs []Transaction
	for _, tx := range txPool.Transactions {
		var outputAmount int
		for _, out := range tx.Outputs {
			outputAmount += out.Amount
		}

		if outputAmount != tx.Input.Amount {
			log.Printf("Invalid sum output amount %d, expected %d\n", outputAmount, tx.Input.Amount)
			continue
		}

		if !VerifyTransaction(&tx) {
			log.Println("Transaction has invalid signature")
			continue
		}

		validTxs = append(validTxs, tx)
	}

	return validTxs
}

func (txPool *TransactionPool) Clear() {
	txPool.Transactions = []Transaction{}
}
