package exchange

import "github.com/AnkurTiwari21/queue"


type ExchangeType string

const ExchangeTypeTopicExchange ExchangeType = "TOPIC"
const ExchangeTypeFanoutExchange ExchangeType = "FANOUT"
const ExchangeTypeDirectExchange ExchangeType = "DIRECT"

type Exchange struct {
	Type   ExchangeType
	Name   string
	Queues map[string]queue.Queue
}

