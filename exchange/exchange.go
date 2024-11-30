package exchange

import (
	"encoding/json"

	"github.com/AnkurTiwari21/binding"
	"github.com/AnkurTiwari21/objects"
	"github.com/AnkurTiwari21/queue"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type ExchangeType string

const ExchangeTypeTopicExchange ExchangeType = "TOPIC"
const ExchangeTypeFanoutExchange ExchangeType = "FANOUT"
const ExchangeTypeDirectExchange ExchangeType = "DIRECT"

type Exchange struct {
	Type   ExchangeType
	Name   string
	Queues map[string]*queue.Queue
}

func (ex *Exchange) SubscribeQueue(write chan []byte, request objects.CommunicationMessage, userConn *websocket.Conn) {
	//check if the queue exixts or not
	q := ex.Queues[request.QueueName]
	logrus.Info("exchange is ", ex.Name)
	logrus.Info("exchange is ", ex.Type)
	logrus.Info("exchange is ", ex.Queues)
	logrus.Info("q is", q)
	if q == nil {
		//no queue is present
		//create a binding
		binding := []binding.Binding{
			{BindingKey: request.BindingKey},
		}
		//create a queue with the name in the request and add the mapping
		q := queue.Queue{
			Name:               request.QueueName,
			Queue:              []string{},
			BindingKey:         binding,
			UserSubscribedConn: []*websocket.Conn{userConn},
		}
		ex.Queues[request.QueueName] = &q
	} else {
		//add the ws connection to the queue
		q.UserSubscribedConn = append(q.UserSubscribedConn, userConn)
	}
	//write a message to write chan to update the user about the operation
	response := objects.Response{
		Message:    "Successfully subscribed to the queue",
		UserWSConnUUID: request.UserConnString,
	}
	logrus.Info("conn")
	logrus.Info(userConn)
	//marhall the response and send to write chan
	convertedResponse, _ := json.Marshal(response)
	write <- convertedResponse
}
