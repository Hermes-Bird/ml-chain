package app

import (
	"fmt"
	"github.com/Hermes-Bird/ml-chain/blockchain"
	"github.com/Hermes-Bird/ml-chain/dto"
	"github.com/Hermes-Bird/ml-chain/p2p"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"strings"
)

func Start() {
	wsPort := os.Getenv("WS_PORT")
	port := os.Getenv("HTTP_PORT")
	envPeers := os.Getenv("PEERS")

	peers := strings.Split(envPeers, ",")

	chain := blockchain.NewBlockchain()

	p2pServer := p2p.NewP2PServer(chain)
	go p2pServer.Start(wsPort)
	p2pServer.ConnectPeers(peers)
	app := fiber.New()

	app.Get("/blocks", func(ctx *fiber.Ctx) error {
		return ctx.JSON(chain.Chain)
	})

	app.Post("/blocks", func(ctx *fiber.Ctx) error {
		var data dto.BlockDataDto
		err := ctx.BodyParser(&data)
		if err != nil {
			ctx.Status(fiber.StatusBadRequest)
			return ctx.JSON(err.Error())
		}

		chain.AddBlock(data.Data)

		return ctx.Redirect("/blocks")
	})

	log.Printf("Starting api-server on port %s", port)
	err := app.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal("Cannot setup http-server", err.Error())
	}
}
