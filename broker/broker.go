package broker

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AnkurTiwari21/exchange"
	"github.com/AnkurTiwari21/objects"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

//broker will have exchanges inside it
//it will run on a port

// {
// 	TASK: SUBSCRIBE_QUEUE
// 	BINDING_KEY:
// 	QUEUE_NAME:
// 	EXCHANGE_NAME:
// 	DURABILITY:
// }

// {
// 	TASK: UNSUBSCRIBE_QUEUE
// 	BINDING_KEY:
// 	QUEUE_NAME:
// 	EXCHANGE_NAME:
// 	DURABILITY:
// }

// {
// 	TASK: PUBLISH_QUEUE
// 	BINDING_KEY:
// 	QUEUE_NAME:
// 	EXCHANGE_NAME:
// }

// {
// 	TASK: CONSUME_QUEUE
// 	QUEUE_NAME:
// 	EXCHANGE_NAME:
// }

//no need of below
// {
// 	TASK: CONFIGURE_EXCHANGE
// 	EXCHANGE_NAME:
// 	EXCHANGE_TYPE:
// 	DURABILITY:
// }

// type MessageToExchange struct {
// 	Task         string          `json:"task,omitempty"`
// 	BindingKey   string          `json:"binding_key,omitempty"`
// 	QueueName    string          `json:"queue_name,omitempty"`
// 	ExchangeType string          `json:"exchange_name,omitempty"`
// 	Durability   bool            `json:"durability,omitempty"`
// 	Status       string          `json:"status,omitempty"`
// 	Message      string          `json:"message,omitempty"`

// }

type Broker struct {
	Port      int32
	Url       string
	Exchanges map[string]exchange.Exchange
	UsersConn map[string]*websocket.Conn
}

// upgrader will help in upgrading the http connection
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for simplicity; adjust as needed
		return true
	},
}

func (br *Broker) InitConnection(read chan []byte, write chan []byte) {
	// Start WebSocket server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleSocketConnection(w, r, br, read, write)
	})

	logrus.Infof("WebSocket server started at ws://%s", br.Url)
	err := http.ListenAndServe(fmt.Sprintf(":%d", br.Port), nil)
	if err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}

func handleSocketConnection(w http.ResponseWriter, r *http.Request, broker *Broker, read chan []byte, write chan []byte) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Errorf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	logrus.Info("New WebSocket connection established!")

	//create a uuid and map that uuid to the connction object so later we can use it
	id := uuid.New()
	uuidStr := id.String()
	broker.UsersConn[uuidStr] = conn

	for {
		// Read message from client
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logrus.Errorf("Error reading message: %v", err)
			break
		}

		var request objects.CommunicationMessage
		err = json.Unmarshal(msg, &request)
		if err != nil {
			logrus.Errorf("Invalid JSON message: %v", err)
			continue
		}

		logrus.Infof("Received task: %s", request.Task)
		//add the connection  uuid
		request.UserConnString = uuidStr
		//convrt to bytes and send
		respBytes, err := json.Marshal(request)
		if err != nil {
			logrus.Error("error marshalling json:err", err)
			return
		}
		//plan: send this task to go routine(read chan)
		read <- respBytes
		//write corresponding result in write chan
	}
	close(read)
	close(write)
}

func (br *Broker) Worker(read chan []byte, write chan []byte, broker *Broker) {
	//here we continuosly accept the data coming and do operations
	for val := range read {
		var request objects.CommunicationMessage
		json.Unmarshal(val, &request)
		logrus.Info("received request from chan : ", request)
		//write logic for processing these messages
		switch request.Task {
		case "SUBSCRIBE_QUEUE":
			exchange := broker.Exchanges[request.ExchangeType]
			exchange.SubscribeQueue(write, request, broker.UsersConn[request.UserConnString])
		case "UNSUBSCRIBE_QUEUE":
			exchange := broker.Exchanges[request.ExchangeType]
			exchange.UnSubscribeQueue(write, request, broker.UsersConn[request.UserConnString])
		case "PUBLISH_QUEUE":
			exchange := broker.Exchanges[request.ExchangeType]
			exchange.Publish(write, request, broker.UsersConn[request.UserConnString])
		case "CONSUME_QUEUE":
			exchange := broker.Exchanges[request.ExchangeType]
			exchange.Consume(write, request, broker.UsersConn[request.UserConnString])
		case "LIST_DATA":
			exchange := broker.Exchanges[request.ExchangeType]
			exchange.ListAllMessages(broker.UsersConn[request.UserConnString])
		}

	}
}
