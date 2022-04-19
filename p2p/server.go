package p2p

import (
	"github.com/Hermes-Bird/ml-chain/blockchain"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)

var upgrader = websocket.Upgrader{}

type P2PServer struct {
	Blockchain *blockchain.Blockchain
	Sockets    []*websocket.Conn
}

func NewP2PServer(blockchain *blockchain.Blockchain) *P2PServer {
	return &P2PServer{
		Blockchain: blockchain,
		Sockets:    []*websocket.Conn{},
	}
}

func (s *P2PServer) Start(port string) {
	server := http.NewServeMux()

	server.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			log.Println("Error while upgrade WS connection", err)
		}
		s.AddSocket(ws)
	})

	log.Printf("Starting p2p-server on port %s", port)
	err := http.ListenAndServe(port, server)
	if err != nil {
		log.Fatalln(err)
	}
}

func (s *P2PServer) ConnectPeers(peers []string) {
	for _, peer := range peers {
		u := url.URL{
			Scheme: "ws",
			Path:   "/",
			Host:   peer,
		}
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Println("Error while connection to peer: ", u.String())
			return
		}
		s.AddSocket(c)
	}
}

func (s *P2PServer) AddSocket(c *websocket.Conn) {
	s.Sockets = append(s.Sockets, c)
	go s.HandleMessages(c)
	err := c.WriteJSON(s.Blockchain.Chain)
	if err != nil {
		log.Println("Error to send blockchain to socket", c.RemoteAddr().String())
	}
	log.Println("Socket connected")
}

func (s *P2PServer) HandleMessages(con *websocket.Conn) {
	for {
		var message []blockchain.Block
		err := con.ReadJSON(&message)
		if err != nil {
			if websocket.IsCloseError(err) {
				return
			}
			log.Println("Error reading message", err.Error())
		}
		log.Println("Message received ->", message)
	}
}
