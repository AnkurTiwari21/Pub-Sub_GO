package exchange

import (
	"encoding/json"
	"math"
	"strings"

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
	if q == nil {
		//no queue is present
		//create a binding
		binding := []string{
			request.BindingKey,
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
		Message:        "Successfully subscribed to the queue",
		UserWSConnUUID: request.UserConnString,
	}
	//marhall the response and send to write chan
	convertedResponse, _ := json.Marshal(response)
	write <- convertedResponse
}

func (ex *Exchange) UnSubscribeQueue(write chan []byte, request objects.CommunicationMessage, userConn *websocket.Conn) {
	//check if the queue exixts or not
	q := ex.Queues[request.QueueName]
	if q != nil {
		//go through all the user connections and remove the connection userConn
		temp := []*websocket.Conn{}
		for _, val := range q.UserSubscribedConn {
			if val != userConn {
				temp = append(temp, val)
			}
		}
		ex.Queues[request.QueueName].UserSubscribedConn = temp
	}
	//write a message to write chan to update the user about the operation
	response := objects.Response{
		Message:        "Successfully unsubscribed to the queue",
		UserWSConnUUID: request.UserConnString,
	}
	//marhall the response and send to write chan
	convertedResponse, _ := json.Marshal(response)
	write <- convertedResponse
}

func (ex *Exchange) Publish(write chan []byte, request objects.CommunicationMessage, userConn *websocket.Conn) {
	//required
	//routingkey
	//message
	//exchangetype

	switch request.ExchangeType {
	case "FANOUT":
		//send the message into all the queues that are present in the exchange --> use workerpool to concurrently update the queue
		distributor := make(chan *queue.Queue)

		go func() {
			for _, val := range ex.Queues {
				distributor <- val
			}
		}()

		go func() {
			//receive a queue and update its content
			q := <-distributor
			q.Enqueue(request.Message)

			//start a seperate go routine to notify all the users connected to the queue that new item is added
			go func() {
				for _, val := range q.UserSubscribedConn {
					response := objects.CommunicationMessage{
						Message: "The queue has a new item",
					}
					resBytes, _ := json.Marshal(&response)
					val.WriteMessage(websocket.TextMessage, resBytes)
				}
			}()
		}()
	case "DIRECT":
		//iterate on queue and whose-ever binding key matches with the routing key the add the message in the queue and update all the conn
		distributor := make(chan *queue.Queue)

		go func() {
			for _, val := range ex.Queues {
				isContain := false
				for _, key := range val.BindingKey {
					if key == request.RoutingKey {
						isContain = true
					}
				}
				if isContain == true {
					distributor <- val
				}
			}
		}()

		go func() {
			//receive a queue and update its content
			q := <-distributor
			q.Enqueue(request.Message)

			//start a seperate go routine to notify all the users connected to the queue that new item is added
			go func() {
				for _, val := range q.UserSubscribedConn {
					response := objects.CommunicationMessage{
						Message: "The queue has a new item",
					}
					resBytes, _ := json.Marshal(&response)
					val.WriteMessage(websocket.TextMessage, resBytes)
				}
			}()
		}()
	case "TOPIC":
		//here the binding key is like ankur.*.uni or #.ozzy.*
		//we are going to implement keys for 3 letters(or seperated by 2 dots)
		// * means 1 word and # means 0 or more words
		//routing key will not contain * or #
		distributor := make(chan *queue.Queue)
		//pick the routing key and split at .
		routingKeySlice := strings.Split(request.RoutingKey, ".")

		//maintain a list of queues where we have to send this message
		queuesList := []*queue.Queue{}

		for _, q := range ex.Queues {
			isContains := true
			for _, key := range q.BindingKey {
				bindingKeySlice := strings.Split(key, ".")
				logrus.Info("binding key slice ", bindingKeySlice)
				for i := 0; i < int(math.Min(float64(len(routingKeySlice)), float64(len(bindingKeySlice)))); i++ {
					if bindingKeySlice[i] == "*" || (bindingKeySlice[i] == routingKeySlice[i]) {
						continue
					}
					if bindingKeySlice[i] == "#" {
						break
					}
					if bindingKeySlice[i] != routingKeySlice[i] {
						isContains = false
						break
					}
				}

				if isContains == true {
					break
				}
			}
			if isContains {
				queuesList = append(queuesList, q)
			}
		}
		logrus.Info("list queue", queuesList)
		logrus.Info("routing key slice", routingKeySlice)
		go func() {
			for _, val := range queuesList {
				distributor <- val
			}
		}()
		go func() {
			//receive a queue and update its content
			q := <-distributor
			q.Enqueue(request.Message)

			//start a seperate go routine to notify all the users connected to the queue that new item is added
			go func() {
				for _, val := range q.UserSubscribedConn {
					response := objects.CommunicationMessage{
						Message: "The queue has a new item",
					}
					resBytes, _ := json.Marshal(&response)
					val.WriteMessage(websocket.TextMessage, resBytes)
				}
			}()
		}()
	}
}
