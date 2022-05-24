package api

import (
	"fmt"
	"github.com/Hermes-Bird/ml-chain/blockchain"
	"github.com/Hermes-Bird/ml-chain/dto"
	"github.com/Hermes-Bird/ml-chain/miner"
	"github.com/Hermes-Bird/ml-chain/p2p"
	"github.com/Hermes-Bird/ml-chain/wallet"
	"github.com/gofiber/fiber/v2"
	"log"
)

func Start(
	port string,
	chain *blockchain.Blockchain,
	p2pServer *p2p.P2PServer,
	wall *wallet.Wallet,
	txPool *wallet.TransactionPool,
	miner *miner.Miner,
) {
	app := fiber.New()

	app.Get("/transactions", func(ctx *fiber.Ctx) error {
		return ctx.JSON(txPool.Transactions)
	})

	app.Post("/transact", func(ctx *fiber.Ctx) error {
		var txData dto.Transact
		err := ctx.BodyParser(&txData)
		if err != nil {
			ctx.Status(fiber.StatusBadRequest)
			return ctx.JSON(err.Error())
		}
		tx := wall.CreateTransaction(txData.Recipient, txData.Amount, chain.GetChainTransactions(), txPool)
		p2pServer.BroadcastTransaction(tx)
		return ctx.Redirect("/transactions")
	})

	app.Get("/blocks", func(ctx *fiber.Ctx) error {
		return ctx.JSON(chain.Chain)
	})
	
	app.Get("/mine-transactions", func(ctx *fiber.Ctx) error {
		block := miner.Mine()
		log.Printf("New block added %s\n", block.Hash)
		return ctx.Redirect("/blocks")
	})

	app.Get("/public-key", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"address": wall.Address,
		})
	})

	log.Printf("Starting api-server on port %s", port)
	err := app.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal("Cannot setup http-server", err.Error())
	}
}
