package kafka_clients

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"google.golang.org/genproto/googleapis/type/date"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/signal"
)

type Instrument struct {
	ShortName      string
	LongName       string
	ISIN           string
	Currency       string
	Market         string
	LotSize        int
	ExpirationDate date.Date
}

var (
	brokerList = kingpin.Flag("brokerList", "List of brokers to connect").Default("localhost:9092").Strings()
	topic      = kingpin.Flag("topic", "Topic name").Default("instruments_topic").String()
	//partition         = kingpin.Flag("partition", "Partition number").Default("0").String()
	//offsetType        = kingpin.Flag("offsetType", "Offset Type (OffsetNewest | OffsetOldest)").Default("-1").Int()
	messageCountStart = kingpin.Flag("messageCountStart", "Message counter start from:").Int()
)

func main() {
	kingpin.Parse()
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
		var instrument Instrument
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
				fmt.Println(instrument.ShortName)
			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
	fmt.Println("Processed", *messageCountStart, "messages")
}
