package main

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/gocql/gocql"
	"github.com/ndjordjevic/go-sb/internal/common"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
	"os/signal"
)

func main() {
	kingpin.Parse()

	// connect to Cassandra cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "go_sb"
	session, _ := cluster.CreateSession()
	log.Println("Connected to Cassandra.")
	defer session.Close()

	// connect to Kafka
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	brokers := *common.BrokerList
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()
	log.Println("Waiting on new messages...")

	instrumentConsumer, err := master.ConsumePartition(*common.InstrumentTopic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	userConsumer, err := master.ConsumePartition(*common.UserTopic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})
	go func() {
		var instrument common.Instrument
		for {
			select {
			case err := <-instrumentConsumer.Errors():
				fmt.Println(err)
			case msg := <-instrumentConsumer.Messages():
				*common.InstrumentMessageCountStart++
				log.Println("Received messages", string(msg.Key), string(msg.Value))
				err := json.Unmarshal(msg.Value, &instrument)
				if err != nil {
					return
				}

				date := instrument.ExpirationDate.Format("2006-01-02")
				insertSql := "INSERT INTO instruments (market, isin, currency, instrument_key, short_name, long_name, expiration_date, status) VALUES ('" + instrument.Market + "', '" + instrument.ISIN + "', '" + instrument.Currency + "', '" + instrument.InstrumentKey + "', '" + instrument.ShortName + "', '" + instrument.LongName + "', '" + date + "', '" + instrument.Status + "')"
				if err := session.Query(insertSql).Exec(); err != nil {
					log.Fatal(err)
				} else {
					log.Println("Instrument is inserted/updated in Cassandra")
				}
			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	go func() {
		var user common.User
		for {
			select {
			case err := <-userConsumer.Errors():
				fmt.Println(err)
			case msg := <-userConsumer.Messages():
				*common.UserMessageCountStart++
				log.Println("Received messages", string(msg.Key), string(msg.Value))
				err := json.Unmarshal(msg.Value, &user)
				if err != nil {
					return
				}

				var accounts = ""
				for index, account := range user.Accounts {
					accTemp := fmt.Sprintf("{currency: '%s', balance: %g}", account.Currency, account.Balance)
					if index < len(user.Accounts)-1 {
						accTemp += ", "
					}
					accounts += accTemp
				}

				userInsertSql := fmt.Sprintf("INSERT INTO users (company, email, first_name, last_name, password, address, city, country, accounts) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', {%s})", user.Company, user.Email, user.FirstName, user.LastName, user.Password, user.Address, user.City, user.Country, accounts)

				if err := session.Query(userInsertSql).Exec(); err != nil {
					log.Fatal(err)
				} else {
					log.Println("User is inserted/updated in Cassandra")
				}
			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	fmt.Println("Processed", *common.InstrumentMessageCountStart, "instrument messages")
	fmt.Println("Processed", *common.UserMessageCountStart, "user messages")
}
