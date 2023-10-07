package producer

import (
	"fmt"

	"github.com/aniketDinda/zocket/models"
	"github.com/streadway/amqp"
)

func Producer(product *models.Product) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer ch.Close()

	_, err = ch.QueueDeclare("ProdImage", false, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	id := product.ProductID

	err = ch.Publish("", "ProdImage", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(id.Hex()),
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
