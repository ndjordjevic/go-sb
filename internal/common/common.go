package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"gopkg.in/alecthomas/kingpin.v2"
	"time"
)

var (
	BrokerList                  = kingpin.Flag("brokerList", "List of brokers to connect").Default("localhost:9092").Strings()
	InstrumentTopic             = kingpin.Flag("instrument-topic", "Instrument topic name").Default("instruments_topic").String()
	UserTopic                   = kingpin.Flag("user-topic", "User topic name").Default("users_topic").String()
	InstrumentMessageCountStart = kingpin.Flag("instrumentMessageCountStart", "Instrument message counter start from:").Int()
	UserMessageCountStart       = kingpin.Flag("userMessageCountStart", "User message counter start from:").Int()
	MaxRetry                    = kingpin.Flag("maxRetry", "Retry limit").Default("5").Int()
)

type Instrument struct {
	Market         string
	ISIN           string
	Currency       string
	InstrumentKey  string
	ShortName      string
	LongName       string
	ExpirationDate time.Time
	Status         string
	Price          float32
}

type Account struct {
	Currency string
	Balance  float64
}

type User struct {
	Company   string
	Email     string
	FirstName string
	LastName  string
	Password  string
	Address   string
	City      string
	Country   string
	Accounts  []Account
}

type Order struct {
	UUID          gocql.UUID
	Email         string
	InstrumentKey string
	Currency      string
	Size          float32
	Price         float32
	Status        string
	Created       time.Time
}

type InstrumentPrice struct {
	InstrumentKey string
	Price         float32
}

func ConvertToByteArray(object interface{}) []byte {
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(object)

	if err != nil {
		fmt.Printf("Error during object encoding: %v", err)
		return nil
	}

	return reqBodyBytes.Bytes()
}
