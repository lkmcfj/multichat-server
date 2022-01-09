package main

import (
	"encoding/json"
	"errors"

	"github.com/gorilla/websocket"
)

type generalPacket struct {
	Action string `json:"action"`
}

type clientRegister struct {
	Action     string `json:"action"`
	ClientName string `json:"client-name"`
	SecretKey  string `json:"secret-key"`
}

func (clientRegister) GetID() int {
	return 1
}

func (clientRegister) GetName() string {
	return "register"
}

func (this *clientRegister) Decode(msg []byte) error {
	err := json.Unmarshal(msg, this)
	if err != nil {
		return err
	}
	if (this.Action != this.GetName()) || (this.ClientName == "") || (this.SecretKey == "") {
		return errors.New("invalid client register body")
	}
	return nil
}

type registerAck struct {
	Action string `json:"action"`
}

func (registerAck) GetID() int {
	return 2
}

func (registerAck) GetName() string {
	return "register-ack"
}

func (this *registerAck) Construct() {
	this.Action = this.GetName()
}

type clientMessage struct {
	Action  string `json:"action"`
	Content string `json:"content"`
}

func (clientMessage) GetID() int {
	return 3
}

func (clientMessage) GetName() string {
	return "client-message"
}

func (this *clientMessage) Decode(msg []byte) error {
	err := json.Unmarshal(msg, this)
	if err != nil {
		return err
	}
	if (this.Action != this.GetName()) || (this.Content == "") {
		return errors.New("invalid client message body")
	}
	return nil
}

type forwardingMessage struct {
	Action           string `json:"action"`
	SourceClientName string `json:"source-client-name"`
	Content          string `json:"content"`
}

func (forwardingMessage) GetID() int {
	return 4
}

func (forwardingMessage) GetName() string {
	return "forwarding-message"
}

func (this *forwardingMessage) Construct(clientName string, content string) {
	this.Action = this.GetName()
	this.SourceClientName = clientName
	this.Content = content
}

func recvPacket(c *websocket.Conn) (int, []byte, error) {
	mt, message, err := c.ReadMessage()
	if err != nil {
		return -1, nil, err
	}
	if mt != websocket.TextMessage {
		return 0, nil, errors.New("wrong websocket message type")
	}
	var getAction generalPacket
	err = json.Unmarshal(message, &getAction)
	if err != nil {
		return 0, message, errors.New("json decode error")
	}
	if (getAction.Action == clientRegister{}.GetName()) {
		return clientRegister{}.GetID(), message, nil
	}
	if (getAction.Action == registerAck{}.GetName()) {
		return registerAck{}.GetID(), message, nil
	}
	if (getAction.Action == clientMessage{}.GetName()) {
		return clientMessage{}.GetID(), message, nil
	}
	if (getAction.Action == forwardingMessage{}.GetName()) {
		return forwardingMessage{}.GetID(), message, nil
	}
	return 0, message, errors.New("unrecognized message type" + getAction.Action)
}
