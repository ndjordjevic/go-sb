package kafka_clients

import "google.golang.org/genproto/googleapis/type/date"

type Instrument struct {
	ShortName      string
	LongName       string
	ISIN           string
	Currency       string
	Market         string
	LotSize        int
	ExpirationDate date.Date
}
