package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/AnkurTiwari21/broker"
	"github.com/AnkurTiwari21/exchange"
	"github.com/AnkurTiwari21/objects"
	"github.com/sirupsen/logrus"
)

func main() {
	//load exchanges(3 predefined)
	mapStringToExchange := map[string]exchange.Exchange{}

	mapStringToExchange["DIRECT-EXCHANGE"] = exchange.Exchange{
		Type: exchange.ExchangeTypeDirectExchange,
		Name: "DIRECT-EXCHANGE",
	}

	mapStringToExchange["FANOUT-EXCHANGE"] = exchange.Exchange{
		Type: exchange.ExchangeTypeFanoutExchange,
		Name: "FANOUT-EXCHANGE",
	}

	mapStringToExchange["TOPIC-EXCHANGE"] = exchange.Exchange{
		Type: exchange.ExchangeTypeTopicExchange,
		Name: "TOPIC-EXCHANGE",
	}

	//init a broker instance and open a tcp conn
	broker := broker.Broker{
		Port:      5672,
		Url:       "127.0.0.1:5672",
		Exchanges: mapStringToExchange,
	}

	listener := broker.InitConnection()

	//on this listener we are going to receive the http request

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

	for {
		conn, err := listener.Accept()
		if err != nil {
			logrus.Info("An error occured while listening...")
			break
		}
		logrus.Info("New connection formed!")
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Read the request headers
	var headers bytes.Buffer
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading request:", err)
			return
		}

		// Append the line to the headers buffer
		headers.WriteString(line)

		// Break if we've reached the end of the headers
		if line == "\r\n" {
			break
		}
	}

	// Extract headers as a string
	headersStr := headers.String()
	fmt.Println("Headers:\n", headersStr)

	// Find the Content-Length header
	contentLength := 0
	for _, line := range strings.Split(headersStr, "\r\n") {
		if strings.HasPrefix(strings.ToLower(line), "content-length:") {
			lengthStr := strings.TrimSpace(strings.Split(line, ":")[1])
			contentLength, _ = strconv.Atoi(lengthStr)
			break
		}
	}

	// Read the body if Content-Length > 0
	if contentLength > 0 {
		body := make([]byte, contentLength)
		_, err := reader.Read(body)
		if err != nil {
			fmt.Println("Error reading body:", err)
			return
		}

		var request objects.HTTPRequestFormat
		// fmt.Println("Body:\n", string(body))
		json.Unmarshal(body, &request)
		logrus.Info("request:")
		logrus.Info("type: ",request.Task)
	}

	// Respond to the client
	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Connection: keep-alive\r\n" + // Keep-alive
		"\r\n" +
		"Received your request!\r\n"

	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing response:", err)
		return
	}
}
