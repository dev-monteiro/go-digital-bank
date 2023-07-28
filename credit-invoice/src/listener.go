package main

import (
	"fmt"

	rabbitmq "github.com/rabbitmq/amqp091-go"
)

type Listener struct {
	Conn *rabbitmq.Connection
	Ch   *rabbitmq.Channel
}

func NewListener() Listener {
	var list Listener
	var err error

	list.Conn, err = rabbitmq.Dial("amqp://guest:guest@rabbitmq")
	if err != nil {
		fmt.Println(err)
	}

	list.Ch, err = list.Conn.Channel()
	if err != nil {
		fmt.Println(err)
	}

	que, err := list.Ch.QueueDeclare("purchases", false, false, false, false, nil)
	msgCh, err := list.Ch.Consume(que.Name, "", true, false, false, false, nil)

	go func() {
		for msg := range msgCh {
			fmt.Println("Message received: " + string(msg.Body))
		}
	}()

	fmt.Println("Listener set up.")
	return list
}

func (lis Listener) Close() {
	lis.Conn.Close()
	lis.Ch.Close()
	fmt.Println("Listener cleaned up.")
}
