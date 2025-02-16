package rabbit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	rabbitmq "github.com/rabbitmq/amqp091-go"
	"log"
	"payment-service/configs"
	"payment-service/internal/entity"
	"time"
)

type PayI interface {
	Update(user *entity.UpdateBalance)
}

type Rabbit struct {
	con         *rabbitmq.Connection
	channel     *rabbitmq.Channel
	store       PayI
	exitContext context.Context
	rabbitCfg   configs.RabbitMQConfig
}

func NewRabbit(ctx context.Context, cfg configs.RabbitMQConfig, store PayI) (*Rabbit, error) {
	newRabbit := &Rabbit{
		store:       store,
		exitContext: ctx,
		rabbitCfg:   cfg,
	}

	err := newRabbit.NewConnection()
	if err != nil {
		return nil, err
	}

	return newRabbit, nil
}

func (myRabbit *Rabbit) Updater(messages <-chan rabbitmq.Delivery) {

	for {
		select {
		case <-myRabbit.exitContext.Done():
			myRabbit.CloseConnection()
			return
		case message := <-messages:
			//TODO:пофикить пустые сообщения
			if len(message.Body) == 0 {
				continue
			}
			var user *entity.UpdateBalance
			err := json.Unmarshal(message.Body, &user)

			if err != nil {
				log.Println(err)
				continue
			}
			myRabbit.store.Update(user)
		}
	}
}

func (myRabbit *Rabbit) NewConnection() error {
	con, err := rabbitmq.Dial(myRabbit.rabbitCfg.RabbitUrl)
	if err != nil {
		return fmt.Errorf("channel is closed: %w", err)
	}

	channel, err := con.Channel()
	if err != nil {
		return fmt.Errorf("Failed to open a channel: %w ", err)
	}

	_, err = channel.QueueDeclare(
		myRabbit.rabbitCfg.NameOfQueue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Failed to declare a queue: %w ", err)
	}

	msgs, err := channel.Consume(
		myRabbit.rabbitCfg.NameOfQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Error on consume message %w ", err)
	}
	if channel.IsClosed() {
		return errors.New("channel is closed")
	}

	log.Println("Connect to RabbitMQ: Success! ")
	go myRabbit.Updater(msgs)
	go myRabbit.Reconnect()

	return nil
}

func (myRabbit *Rabbit) Reconnect() {
	closeErrChan := make(chan *rabbitmq.Error)
	if myRabbit.con != nil {
		myRabbit.con.NotifyClose(closeErrChan)
	}
	for {
		<-closeErrChan
		log.Println("Trying to Reconnect...")
		timer := 1
		for {
			time.Sleep(time.Duration(timer) * time.Second)
			err := myRabbit.NewConnection()
			if err != nil {
				timer *= 2
				if timer > 30 {
					timer = 30
				}
				log.Println("Failed to reconnect to RabbitMQ. Retrying...")
				continue
			}
			log.Print("Reconnected to RabbitMQ")
			break
		}

	}
}

func (myRabbit *Rabbit) CloseConnection() {
	err := myRabbit.con.Close()
	if err != nil {
		log.Println("Failed to close connection")
	}
	err = myRabbit.channel.Close()
	if err != nil {
		log.Println("Failed to close channel")
	}
	log.Println("Connection and channel are closed")
}
