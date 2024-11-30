package queue

import (
	"github.com/gorilla/websocket"
)

type Queue struct {
	Name               string
	Queue              []string //queue DS
	BindingKey         []string
	UserSubscribedConn []*websocket.Conn
}

func (q *Queue) Enqueue(str string) {
	q.Queue = append(q.Queue, str)
}

func (q *Queue) Dequeue() string {
	if len(q.Queue) > 0 {
		queueFront := q.Queue[0]
		temp := q.Queue[1:]
		q.Queue = temp
		return queueFront
	} else {
		return "NO ELEMENT PRESENT"
	}
}
