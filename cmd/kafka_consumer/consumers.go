package main

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/gocql/gocql"
	"github.com/ndjordjevic/go-sb/internal/kafka_common"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
	"os/signal"
	"strconv"
)

var (
	brokerList = kingpin.Flag("brokerList", "List of brokers to connect").Default("localhost:9092").Strings()
	topic      = kingpin.Flag("topic", "Topic name").Default("instruments_topic").String()
	//partition         = kingpin.Flag("partition", "Partition number").Default("0").String()
	//offsetType        = kingpin.Flag("offsetType", "Offset Type (OffsetNewest | OffsetOldest)").Default("-1").Int()
	messageCountStart = kingpin.Flag("messageCountStart", "Message counter start from:").Int()
)

func main() {
	kingpin.Parse()

	// connect to Cassandra the cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "go_sb"
	session, _ := cluster.CreateSession()
	defer session.Close()

	// connect to Kafka
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	brokers := *brokerList
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()
	consumer, err := master.ConsumePartition(*topic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})
	go func() {
		var instrument kafka_common.Instrument
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				*messageCountStart++
				fmt.Println("Received messages", string(msg.Key), string(msg.Value))
				err := json.Unmarshal(msg.Value, &instrument)
				if err != nil {
					return
				}

				date := strconv.Itoa(int(instrument.ExpirationDate.Year)) + "-" + strconv.Itoa(int(instrument.ExpirationDate.Month)) + "-" + strconv.Itoa(int(instrument.ExpirationDate.Day))
				insertSql := "INSERT INTO instruments (market, isin, currency, short_name, long_name, expiration_date, status) VALUES ('" + instrument.Market + "', '" + instrument.ISIN + "', '" + instrument.Currency + "', '" + instrument.ShortName + "', '" + instrument.LongName + "', '" + date + "', '" + instrument.Status + "')"
				log.Println(insertSql)
				if err := session.Query(insertSql).Exec(); err != nil {
					log.Fatal(err)
				}
			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
	fmt.Println("Processed", *messageCountStart, "messages")
}
