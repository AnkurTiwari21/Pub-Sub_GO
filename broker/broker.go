package broker

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AnkurTiwari21/exchange"
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

// {
// 	TASK: CONFIGURE_EXCHANGE
// 	EXCHANGE_NAME:
// 	EXCHANGE_TYPE:
// 	DURABILITY:
// }

// type RequestMessage struct {
// 	Task         string `json:"task"`
// 	BindingKey   string `json:"binding_key"`
// 	QueueName    string `json:"queue_name"`
// 	ExchangeName string `json:"exchange_name"`
// 	Durability   bool   `json:"durability"`
// }

type CommunicationMessage struct {
	Task         string `json:"task,omitempty"`
	BindingKey   string `json:"binding_key,omitempty"`
	QueueName    string `json:"queue_name,omitempty"`
	ExchangeType string `json:"exchange_name,omitempty"`
	Durability   bool   `json:"durability,omitempty"`
	Status       string `json:"status,omitempty"`
	Message      string `json:"message,omitempty"`
}

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

func (br *Broker) InitConnection() {
	// Start WebSocket server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleSocketConnection(w, r, br)
	})

	logrus.Infof("WebSocket server started at ws://%s", br.Url)
	err := http.ListenAndServe(fmt.Sprintf(":%d", br.Port), nil)
	if err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}

func handleSocketConnection(w http.ResponseWriter, r *http.Request, broker *Broker) {
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

		var request CommunicationMessage
		err = json.Unmarshal(msg, &request)
		if err != nil {
			logrus.Errorf("Invalid JSON message: %v", err)
			continue
		}

		logrus.Infof("Received task: %s", request.Task)

		// exchange:=broker.Exchanges[request.ExchangeType]

		//plan: send this task to go routine(read chan)
		//write corresponding result in write chan

		// Respond to the client
		response := CommunicationMessage{
			Status:  "success",
			Message: fmt.Sprintf("Processed task: %s", request.Task),
		}
		respBytes, _ := json.Marshal(response)
		logrus.Info("connection is ", conn)
		err = conn.WriteMessage(websocket.TextMessage, respBytes)
		if err != nil {
			logrus.Errorf("Error writing response: %v", err)
			break
		}
	}
}
