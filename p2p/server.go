package p2p

import (
	"encoding/json"
	"github.com/Hermes-Bird/ml-chain/blockchain"
	"github.com/Hermes-Bird/ml-chain/dto"
	"github.com/Hermes-Bird/ml-chain/wallet"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)

var upgrader = websocket.Upgrader{}

const TransactionDtoType = "transaction"
const ChainDtoType = "chain"
const ClearDtoType = "clear"

type P2PServer struct {
	Blockchain *blockchain.Blockchain
	Sockets    []*websocket.Conn
	TxPool     *wallet.TransactionPool
}

func NewP2PServer(blockchain *blockchain.Blockchain, txPool *wallet.TransactionPool) *P2PServer {
	return &P2PServer{
		Blockchain: blockchain,
		TxPool:     txPool,
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
	if len(peers) == 0 {
		return
	}
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

	SendMessage(c, &dto.SocketMessage{
		MessageType: ChainDtoType,
		Data:        s.Blockchain.Chain,
	})

	log.Println("Socket connected")
}

func (s P2PServer) SyncChains() {
	for _, ws := range s.Sockets {
		SendMessage(ws, &dto.SocketMessage{
			MessageType: ChainDtoType,
			Data:        s.Blockchain.Chain,
		})
	}
}

func SendMessage(ws *websocket.Conn, dto *dto.SocketMessage) {
	err := ws.WriteJSON(dto)
	if err != nil {
		log.Println("Error to send message to socket", ws.RemoteAddr().String())
	}
}

func (s P2PServer) BroadcastTransaction(tx *wallet.Transaction) {
	for _, ws := range s.Sockets {
		SendMessage(ws, &dto.SocketMessage{
			MessageType: TransactionDtoType,
			Data:        tx,
		})
	}
}

func (s P2PServer) BroadcastClear() {
	for _, ws := range s.Sockets {
		SendMessage(ws, &dto.SocketMessage{
			MessageType: ClearDtoType,
		})
	}
}

func (s *P2PServer) HandleMessages(con *websocket.Conn) {
	for {
		var message dto.SocketMessage
		err := con.ReadJSON(&message)

		if err != nil {
			if websocket.IsCloseError(err) {
				return
			}
			log.Println("Error reading message", err.Error())
		}

		switch message.MessageType {
		case ChainDtoType:
			chain, err := UpcastToType[[]blockchain.Block](message.Data)
			if err == nil {
				s.Blockchain.ReplaceChain(chain)
			}
		case TransactionDtoType:
			tx, err := UpcastToType[wallet.Transaction](message.Data)
			if err == nil {
				s.TxPool.AddOrUpdateTx(&tx)
			}
		case ClearDtoType:
			s.TxPool.Clear()
		}

	}
}

func UpcastToType[T any](data any) (T, error) {
	var res T

	bytes, err := json.Marshal(data)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}
