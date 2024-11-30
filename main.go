package main

import (
	"encoding/json"

	"github.com/AnkurTiwari21/broker"
	"github.com/AnkurTiwari21/exchange"
	"github.com/AnkurTiwari21/models"
	"github.com/AnkurTiwari21/objects"
	"github.com/AnkurTiwari21/queue"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	//declare a read and write channel here
	var read = make(chan []byte)
	var write = make(chan []byte)

	// Initialize exchanges
	mapStringToExchange := map[string]exchange.Exchange{
		"DIRECT": {
			Type:   exchange.ExchangeTypeDirectExchange,
			Name:   "DIRECT",
			Queues: map[string]*queue.Queue{},
		},
		"FANOUT": {
			Type:   exchange.ExchangeTypeFanoutExchange,
			Name:   "FANOUT",
			Queues: map[string]*queue.Queue{},
		},
		"TOPIC": {
			Type:   exchange.ExchangeTypeTopicExchange,
			Name:   "TOPIC",
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

	//recover data that is saved in the db

	//connection
	dsn := "host=localhost user=ankur password=ankur dbname=pubsub port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Error("Error in Establishing Connection :", err)
		return
	}
	db.SetupJoinTable(&models.Exchange{}, "Queues", &models.Binding{})
	db.AutoMigrate(&models.Exchange{}, &models.Queue{}, &models.Binding{})
	//

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
