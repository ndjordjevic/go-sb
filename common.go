package kafka_clients

import "google.golang.org/genproto/googleapis/type/date"

type Instrument struct {
	Market         string
	ISIN           string
	Currency       string
	ShortName      string
	LongName       string
	LotSize        int
	ExpirationDate date.Date
	Status         string
}
