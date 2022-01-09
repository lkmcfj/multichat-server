package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type client struct {
	Name       string
	Valid      bool
	Connection *websocket.Conn
}

var clients []client
var lock sync.Mutex // protects adding/removing client and writing to clients

func addClient(name string, connection *websocket.Conn) int {
	lock.Lock()
	defer lock.Unlock()
	for i := 0; i < len(clients); i++ {
		if !clients[i].Valid {
			clients[i].Valid = true
			clients[i].Name = name
			clients[i].Connection = connection
			return i
		}
	}
	clients = append(clients, client{Name: name, Valid: true, Connection: connection})
	return len(clients) - 1
}

func removeClient(id int) {
	lock.Lock()
	clients[id].Valid = false
	lock.Unlock()
}

func forwardMessage(clientID int, clientName string, content string) {
	var forward forwardingMessage
	forward.Construct(clientName, content)
	msg, err := json.Marshal(forward)
	if err != nil {
		log.Println("error when encoding forwarding message:", err)
		return
	}
	lock.Lock()
	defer lock.Unlock()
	for i := 0; i < len(clients); i++ {
		if i == clientID {
			continue
		}
		if clients[i].Valid {
			err := clients[i].Connection.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("error when forwarding message:", err)
			}
		}
	}
}

func sendRegisterAck(connection *websocket.Conn) error {
	var ack registerAck
	ack.Construct()
	msg, err := json.Marshal(ack)
	if err != nil {
		return err
	}
	return connection.WriteMessage(websocket.TextMessage, msg)
}

func serve(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket upgrade:", err)
		return
	}
	defer c.Close()
	mt, message, err := recvPacket(c)
	if err != nil {
		log.Println("invalid received packet:", err)
		return
	}
	if (mt != clientRegister{}.GetID()) {
		log.Println("the first packet should be client register")
		return
	}
	var regInfo clientRegister
	err = regInfo.Decode(message)
	if err != nil {
		log.Println(err)
		return
	}
	if regInfo.SecretKey != globalConfig.SecretKey {
		log.Println("client register with wrong secret key: refused")
		return
	}
	err = sendRegisterAck(c)
	if err != nil {
		log.Println("error when sending register ack:", err)
		return
	}
	curID := addClient(regInfo.ClientName, c)
	for {
		mt, message, err := recvPacket(c)
		if mt < 0 {
			log.Println("recvPacket fatal error:", err, ", close connection")
			break
		}
		if err != nil {
			log.Println("recvPacket:", err)
			continue
		}
		if (mt != clientMessage{}.GetID()) {
			log.Println("unsupported action type:", mt)
			continue
		}
		var msgInfo clientMessage
		err = msgInfo.Decode(message)
		if err != nil {
			log.Println(err)
			continue
		}
		forwardMessage(curID, regInfo.ClientName, msgInfo.Content)
	}
	removeClient(curID)
}

func main() {
	err := loadConfig()
	if err != nil {
		log.Fatal("fail to load configuration:", err)
		return
	}
	http.HandleFunc(globalConfig.WSpath, serve)
	log.Fatal(http.ListenAndServe(globalConfig.Host, nil))
}
