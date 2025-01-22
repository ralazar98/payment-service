package rabbit

import (
	"encoding/json"
	rab "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"payment-service/internal/entity"
)

type PayI interface {
	Update(user *entity.UpdateBalance)
}

type rabbit struct {
	con         *rab.Connection
	channel     *rab.Channel
	nameOfQueue string
	store       PayI
	rabbitUrl   string
}

func NewRabbit(store PayI) *rabbit {
	nameOfQueue := os.Getenv("NAME_OF_QUEUE")
	rabbitUrl := os.Getenv("RABBIT_URL")

	return &rabbit{
		nameOfQueue: nameOfQueue,
		store:       store,
		rabbitUrl:   rabbitUrl,
	}
}

func (r *rabbit) Updater() {

	msgs, err := r.channel.Consume(
		r.nameOfQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Error on consume message", err)
	}
	if r.channel.IsClosed() {
		log.Println("Channel is closed")
	}

	for message := range msgs {
		var user *entity.UpdateBalance
		err := json.Unmarshal(message.Body, &user)
		if err != nil {
			log.Println(err)
		}
		r.store.Update(user)
	}
}

func (newRabbit *rabbit) NewConnection() {
	con, err := rab.Dial(newRabbit.rabbitUrl)
	if err != nil {
		log.Println(err, "Failed to connect to RabbitMQ")
	}
	newRabbit.con = con
	log.Println("Connected to RabbitMQ")

	channel, err := con.Channel()
	if err != nil {
		log.Println(err, "Failed to open a channel")
	}
	newRabbit.channel = channel

	_, err = channel.QueueDeclare(
		newRabbit.nameOfQueue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err, "Failed to declare a queue")
	}
}

func (newRabbit *rabbit) CloseConnection() {
	err := newRabbit.con.Close()
	if err != nil {
		log.Println("Failed to close connection")
	}
	err = newRabbit.channel.Close()
	if err != nil {
		log.Println("Failed to close channel")
	}
}
