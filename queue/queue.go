package queue

import (
	"github.com/AnkurTiwari21/binding"
)

type Queue struct {
	Name       string
	Queue      []string
	BindingKey []binding.Binding
}

func (q *Queue) Enqueue(str string) {
	q.Queue = append(q.Queue, str)
}

func (q *Queue) Dequeue() string {
	queueFront := q.Queue[0]
	temp := q.Queue[1:]
	q.Queue = temp
	return queueFront
}
