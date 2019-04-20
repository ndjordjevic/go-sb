package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/ndjordjevic/go-sb/internal/kafka_common"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
)

func main() {
	userToSend := kafka_common.User{
		Company:   "FIS",
		Email:     "ggg@gmail.com",
		FirstName: "Gaga",
		LastName:  "Dragic",
		Password:  "gaga123",
		Address:   "Save Simica 27",
		City:      "Smederevo",
		Country:   "Serbia",
		Accounts: []kafka_common.Account{
			{Currency: "USD", Balance: 3000},
			{Currency: "EUR", Balance: 1000},
		},
	}

	byteArray := kafka_common.ConvertToByteArray(userToSend)

	log.Println(byteArray)
	log.Println(string(byteArray))

	kingpin.Parse()
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = *kafka_common.MaxRetry
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(*kafka_common.BrokerList, config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			panic(err)
		}
	}()
	msg := &sarama.ProducerMessage{
		Topic: *kafka_common.UserTopic,
		Value: sarama.ByteEncoder(byteArray),
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", *kafka_common.UserTopic, partition, offset)
}
