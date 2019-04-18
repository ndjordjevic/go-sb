package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/ndjordjevic/go-sb/internal/kafka_common"
	"google.golang.org/genproto/googleapis/type/date"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	brokerList = kingpin.Flag("brokerList", "List of brokers to connect").Default("localhost:9092").Strings()
	topic      = kingpin.Flag("topic", "Topic name").Default("instruments_topic").String()
	maxRetry   = kingpin.Flag("maxRetry", "Retry limit").Default("5").Int()
)

func main() {
	instrumentToSend := kafka_common.Instrument{
		Market:    "Xetra",
		ISIN:      "BMW001",
		Currency:  "SEK",
		ShortName: "BMW",
		LongName:  "BMW Incorporation",
		LotSize:   1,
		ExpirationDate: date.Date{
			Year:  2019,
			Month: 10,
			Day:   31,
		},
		Status: "ACTIVE",
	}

	byteArray := convertToByteArray(instrumentToSend)

	fmt.Println(byteArray)
	fmt.Println(string(byteArray))

	kingpin.Parse()
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = *maxRetry
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(*brokerList, config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			panic(err)
		}
	}()
	msg := &sarama.ProducerMessage{
		Topic: *topic,
		Value: sarama.ByteEncoder(byteArray),
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", *topic, partition, offset)
}

func convertToByteArray(instrument kafka_common.Instrument) []byte {
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(instrument)

	if err != nil {
		fmt.Printf("Error during instrument encoding: %v", err)
		return nil
	}

	return reqBodyBytes.Bytes()
}
