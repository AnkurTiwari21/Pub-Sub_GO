package main

import (
	"github.com/AnkurTiwari21/broker"
	"github.com/AnkurTiwari21/exchange"
	"github.com/gorilla/websocket"
)

func main() {
	//declare a read and write channel here
	
	// Initialize exchanges
	mapStringToExchange := map[string]exchange.Exchange{
		"DIRECT-EXCHANGE": {
			Type: exchange.ExchangeTypeDirectExchange,
			Name: "DIRECT-EXCHANGE",
		},
		"FANOUT-EXCHANGE": {
			Type: exchange.ExchangeTypeFanoutExchange,
			Name: "FANOUT-EXCHANGE",
		},
		"TOPIC-EXCHANGE": {
			Type: exchange.ExchangeTypeTopicExchange,
			Name: "TOPIC-EXCHANGE",
		},
	}

	// Init a broker instance
	broker := broker.Broker{
		Port:      5672,
		Url:       "127.0.0.1:5672",
		Exchanges: mapStringToExchange,
		UsersConn: map[string]*websocket.Conn{},
	}

	broker.InitConnection()
}
