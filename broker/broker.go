package broker

import (
	"net"

	"github.com/AnkurTiwari21/exchange"
	"github.com/sirupsen/logrus"
)

//broker will have exchanges inside it
//it will run on a port

type Broker struct {
	Port      int32
	Url       string
	Exchanges map[string]exchange.Exchange
}

func (br *Broker) InitConnection() net.Listener {
	//open a tcp socket on port 5672
	listener, err := net.Listen("tcp", "127.0.0.1:5672")
	if err != nil {
		logrus.Info("TCP connection started at 127.0.0.1:5672!")
		logrus.Info("LISTENING...")
	}
	return listener
}
