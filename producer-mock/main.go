package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	rabbitmq "github.com/rabbitmq/amqp091-go"
)

type PurchaseEvent struct {
	PurchaseId          int
	CreditAccountId     int
	PurchaseDateTime    string
	Amount              float32
	NumInstallments     int
	MerchantDescription string
	Status              string
}

func main() {
	time.Sleep(10 * time.Second)

	conn, err := rabbitmq.Dial("amqp://guest:guest@rabbitmq")
	checkError(err)
	defer conn.Close()

	cha, err := conn.Channel()
	checkError(err)
	defer cha.Close()

	que, err := cha.QueueDeclare(
		"purchases",
		false,
		false,
		false,
		false,
		nil,
	)
	checkError(err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	time.Sleep(3 * time.Second)
	for i := 1; i <= 60; i++ {
		time.Sleep(3 * time.Second)

		event := PurchaseEvent{
			PurchaseId:          i,
			CreditAccountId:     123,
			PurchaseDateTime:    time.Now().String(),
			Amount:              float32(r.Intn(10000) / 100),
			NumInstallments:     1,
			MerchantDescription: "Acme Shop",
			Status:              "APPROVED",
		}
		body, err := json.Marshal(event)
		checkError(err)

		err = cha.PublishWithContext(
			ctx,
			"",
			que.Name,
			false,
			false,
			rabbitmq.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
		checkError(err)
		fmt.Println("Sent: " + string(body))
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
