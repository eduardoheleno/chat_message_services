package main

import (
	"chat_service/internal/domain/message"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type WsMessage struct {
	SenderId uint `json:"sender_id"`
	ReceiverId uint `json:"receiver_id"`
	ChatId uint `json:"chat_id"`
	ReceiverEmail string `json:"receiver_email"`

	TargetId uint `json:"target_id"`
	Type string `json:"type"`

	Nonce []byte `json:"nonce"`
	Content []byte `json:"content"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func initDatabase() *gorm.DB {
	dbPsswPath := os.Getenv("MYSQL_ROOT_PASSWORD_FILE")
	dbPswd, fileErr := os.ReadFile(dbPsswPath)
	if fileErr != nil {
		log.Panicf("Password file not found: %s", fileErr)
	}

	dsn := fmt.Sprintf("root:%s@tcp(chat_database:3306)/chat_api?charset=utf8&parseTime=true", dbPswd)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("Failed to connect to database: %s", err)
	}

	return db
}

func main() {
	db := initDatabase()
	messageRepository := message.NewMessageRepository(db)

	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	publisherCh, err := conn.Channel()
	failOnError(err, "Failed to open the publisher channel")
	defer publisherCh.Close()

	q, err := ch.QueueDeclare(
		"messages",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	redis := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		Password: "",
		DB: 0,
	})

	var forever chan struct{}

	go func() {
		for d := range msgs {
			var wsMessage WsMessage
			parseErr := json.Unmarshal(d.Body, &wsMessage)
			if parseErr != nil {
				d.Nack(false, true)
				return
			}

			savedMessage := message.Message{
				UserId: wsMessage.SenderId,
				ChatId: wsMessage.ChatId,
				Content: wsMessage.Content,
				Nonce: wsMessage.Nonce,
			}
			repoErr := messageRepository.Create(&savedMessage)
			if repoErr != nil {
				log.Println("Failed to save message: ", repoErr)
				d.Nack(false, true)
				return
			}
			// savedMessage.TargetId = wsMessage.TargetId
			// savedMessage.Type = wsMessage.Type

			userKey := strconv.FormatUint(uint64(wsMessage.ReceiverId), 10)
			redisCtx := context.Background()
			nodeHash, err := redis.Get(redisCtx, userKey).Result()
			if err != nil {
				d.Ack(false)
				continue
			}

			message := map[string]interface{}{
				"target_id": wsMessage.TargetId,
				"type": wsMessage.Type,
				"message": savedMessage,
			}
			jsonMessage, parseErr := json.Marshal(message)

			// TODO: create message hash to allow to Nack back to the queue
			// without storing it again
			if parseErr != nil {
				return
			}

			// TODO: create message hash to allow to Nack back to the queue
			// without storing it again
			publisherCh.Publish(
				"",
				nodeHash,
				false,
				false,
				amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					ContentType: "application/json",
					Body: jsonMessage,
				},
			)

			d.Ack(false)
		}
	}()
	log.Printf("[*] Waiting for messages...")
	<-forever
}
