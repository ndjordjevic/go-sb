package main

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/gomodule/redigo/redis"
	pricepb "github.com/ndjordjevic/go-sb/api/price"
	"google.golang.org/grpc"
	"log"
	"net"
)

type priceServer struct{}

var pool = newPool()

func (priceServer) RequestPrices(context.Context, *empty.Empty) (*pricepb.GetAllPricesResponse, error) {
	// new Redis pool and connection
	conn := pool.Get()
	defer func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()

	var prices []struct {
		InstrumentKey string
		Price         float32
	}

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
