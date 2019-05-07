package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/gomodule/redigo/redis"
	pricepb "github.com/ndjordjevic/go-sb/api/price"
	"github.com/ndjordjevic/go-sb/internal/common"
	"google.golang.org/grpc"
	"log"
	"net"
)

type priceServer struct{}

var pool = newPool()

func (priceServer) StreamPriceChange(req *empty.Empty, stream pricepb.PriceService_StreamPriceChangeServer) error {
	conn := pool.Get()
	defer func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()

	psc := redis.PubSubConn{Conn: conn}
	if err := psc.Subscribe("price_updates"); err != nil {
		log.Fatal(err)
	}

	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			fmt.Printf("%s: message: %s\n", v.Channel, v.Data)

			instrumentPrice := FromGOB64(string(v.Data))
			log.Println(instrumentPrice)

			res := &pricepb.StreamPriceChangeResponse{
				Price: &pricepb.Price{
					InstrumentKey: instrumentPrice.InstrumentKey,
					Price:         instrumentPrice.Price,
				},
			}

			if err := stream.Send(res); err != nil {
				log.Fatal(err)
			}
		case redis.Subscription:
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			return v
		}
	}
}

func (priceServer) GetAllPrices(context.Context, *empty.Empty) (*pricepb.GetAllPricesResponse, error) {
	conn := pool.Get()
	defer func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()

	var prices []common.InstrumentPrice

	values, err := redis.Values(conn.Do("HGETALL", "instrument_prices"))
	if err != nil {
		log.Println(err)
	}

	if err := redis.ScanSlice(values, &prices); err != nil {
		log.Panic(err)
	}

	log.Println(prices)

	var pbprices []*pricepb.Price

	for _, price := range prices {
		pbprices = append(pbprices, &pricepb.Price{
			InstrumentKey: price.InstrumentKey,
			Price:         price.Price,
		})
	}

	res := &pricepb.GetAllPricesResponse{
		Prices: pbprices,
	}

	return res, nil
}

func newPool() *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

// go binary decoder
func FromGOB64(str string) common.InstrumentPrice {
	instrumentPrice := common.InstrumentPrice{}
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		fmt.Println(`failed base64 Decode`, err)
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(&instrumentPrice)
	if err != nil {
		fmt.Println(`failed gob Decode`, err)
	}
	return instrumentPrice
}

func main() {
	// price grpc service server
	lis, err := net.Listen("tcp", "localhost:50071")

	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	pricepb.RegisterPriceServiceServer(s, &priceServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
