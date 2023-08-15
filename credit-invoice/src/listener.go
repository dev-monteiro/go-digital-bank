package main

import (
	"encoding/json"
	"fmt"

	rabbitmq "github.com/rabbitmq/amqp091-go"
)

type Listener struct {
	conn  *rabbitmq.Connection
	chann *rabbitmq.Channel
	repo  PurchaseRepo
}

func NewListener(repo PurchaseRepo) Listener {
	conn, err := rabbitmq.Dial("amqp://guest:guest@rabbitmq")
	if err != nil {
		fmt.Println(err)
	}

	chann, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}

	que, err := chann.QueueDeclare("purchases", false, false, false, false, nil)
	msgCh, err := chann.Consume(que.Name, "", true, false, false, false, nil)

	go func() {
		for msg := range msgCh {
			purchase := Purchase{}
			json.Unmarshal(msg.Body, &purchase)
			repo.save(purchase)
		}
	}()

	fmt.Println("Listener set up.")
	return Listener{
		conn:  conn,
		chann: chann,
		repo:  repo,
	}
}

func (lis *Listener) Close() {
	lis.conn.Close()
	lis.chann.Close()
	fmt.Println("Listener cleaned up.")
}
