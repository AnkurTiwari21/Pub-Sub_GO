package queue

import (
	"github.com/AnkurTiwari21/binding"
	"github.com/AnkurTiwari21/user"
)

type Queue struct {
	Name       string
	Queue      []string
	BindingKey []binding.Binding
	Users      []user.User
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