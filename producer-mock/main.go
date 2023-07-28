package main

import (
	"context"
	"encoding/json"
	"fmt"
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

	for i := 1; i <= 20; i++ {
		time.Sleep(1 * time.Second)

		event := PurchaseEvent{
			PurchaseId:          i,
			CreditAccountId:     123,
			PurchaseDateTime:    "2023-07-01T09:00:00",
			Amount:              19.99,
			NumInstallments:     1,
			MerchantDescription: "Starbucks",
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
