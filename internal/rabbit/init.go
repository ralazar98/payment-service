package rabbit

import (
	"context"
	"encoding/json"
	"errors"
	rabbitmq "github.com/rabbitmq/amqp091-go"
	"log"
	"payment-service/configs"
	"payment-service/internal/entity"
)

type PayI interface {
	Update(user *entity.UpdateBalance)
}

type Rabbit struct {
	con         *rabbitmq.Connection
	channel     *rabbitmq.Channel
	store       PayI
	exitContext context.Context
}

func NewRabbit(ctx context.Context, store PayI) *Rabbit {
	return &Rabbit{
		store:       store,
		exitContext: ctx,
	}
}

func (myRabbit *Rabbit) Updater(messages <-chan rabbitmq.Delivery) {

	for {
		select {
		case <-myRabbit.exitContext.Done():
			log.Println("Test")
			//TODO: убрать лог
			myRabbit.CloseConnection()
			return
		case message := <-messages:
			var user *entity.UpdateBalance
			err := json.Unmarshal(message.Body, &user)
			if err != nil {
				//TODO:обработать ошибку, у тебя есть ошибка, то не должны вызывать следующую функцию
				log.Println(err)
			}
			myRabbit.store.Update(user)
		}
	}
}

func (myRabbit *Rabbit) NewConnection(cfg configs.RabbitMQConfig) error {
	con, err := rabbitmq.Dial(cfg.RabbitUrl)
	if err != nil {
		log.Println("Failed to connect to RabbitMQ :", err)
		return err
	} else {
		log.Println("Connected to RabbitMQ")
	}
	//TODO:Имеет смысл не делать блок else
	// а просто написать так
	// if err != nil {
	//		log.Println("Failed to connect to RabbitMQ :", err)
	//		return err
	//	}
	//	log.Println("Connected to RabbitMQ")
	channel, err := con.Channel()
	if err != nil {
		log.Println("Failed to open a channel:", err)
		return err
	} else {
		log.Println("Channel opened")
	}

	_, err = channel.QueueDeclare(
		cfg.NameOfQueue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Failed to declare a queue:", err)
		return err
	} else {
		log.Println("Queue is declared")
		//TODO:Аналогично с прошлой
	}
	msgs, err := channel.Consume(
		cfg.NameOfQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Error on consume message", err)
		return err
	}
	if channel.IsClosed() {
		log.Println("Channel is closed")
		return errors.New("channel is closed")
	}
	//TODO:По ошибкам
	// подумай насчет того, что бы просто возвращать новую ошибку с тексто, как у тебя в логах, но сохраняя старую
	// return fmt.Errorf("channel is closed: %w", err)
	// а сам лог об ошибки уже вызывать из места где ты его вызываешь, в нашем случае мейн

	go myRabbit.Updater(msgs)

	return nil
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
