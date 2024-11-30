package main

import (
	"encoding/json"

	"github.com/AnkurTiwari21/broker"
	"github.com/AnkurTiwari21/exchange"
	"github.com/AnkurTiwari21/objects"
	"github.com/AnkurTiwari21/queue"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func main() {
	//declare a read and write channel here
	var read = make(chan []byte)
	var write = make(chan []byte)

	// Initialize exchanges
	mapStringToExchange := map[string]exchange.Exchange{
		"DIRECT": {
			Type:   exchange.ExchangeTypeDirectExchange,
			Name:   "DIRECT-EXCHANGE",
			Queues: map[string]*queue.Queue{},
		},
		"FANOUT": {
			Type:   exchange.ExchangeTypeFanoutExchange,
			Name:   "FANOUT-EXCHANGE",
			Queues: map[string]*queue.Queue{},
		},
		"TOPIC": {
			Type:   exchange.ExchangeTypeTopicExchange,
			Name:   "TOPIC-EXCHANGE",
			Queues: map[string]*queue.Queue{},
		},
	}

	// Init a broker instance
	broker := broker.Broker{
		Port:      5672,
		Url:       "127.0.0.1:5672",
		Exchanges: mapStringToExchange,
		UsersConn: map[string]*websocket.Conn{},
	}

	//fire workers to work on read and write chan that will continuously be update in the conn below
	go broker.Worker(read, write, &broker)
	go func() {
		for val := range write {
			var responseObject objects.Response
			err := json.Unmarshal(val, &responseObject)
			if err != nil {
				logrus.Error("error in unmarshalling the response!")
				return
			}
			logrus.Info("response obj ", responseObject)
			logrus.Info(responseObject.UserWSConnUUID)
			broker.UsersConn[responseObject.UserWSConnUUID].WriteMessage(websocket.TextMessage, val)
		}
	}()

	broker.InitConnection(read, write)

}
