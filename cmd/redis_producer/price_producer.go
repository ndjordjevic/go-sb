package main

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/gomodule/redigo/redis"
	"github.com/ndjordjevic/go-sb/internal/common"
	"log"
	"math/rand"
	"time"
)

func main() {
	// connect to Cassandra cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "go_sb"
	session, _ := cluster.CreateSession()
	log.Println("Connected to Cassandra.")
	defer session.Close()

	var instrumentKey string
	// new Redis pool and connection
	pool := newPool()
	conn := pool.Get()
	defer func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()

	for {
		iter := session.Query("SELECT instrument_key FROM instruments").Iter()
		for iter.Scan(&instrumentKey) {
			price := rand.Intn(100)

			_, err := conn.Do("HSET", "instrument_prices", instrumentKey, price)
			if err != nil {
				log.Println(err)
			}

			log.Println("Instrument:", instrumentKey, "Price:", price)

			_, err = conn.Do("PUBLISH", "price_updates", ToGOB64(common.InstrumentPrice{
				instrumentKey, float32(price),
			}))
			if err != nil {
				log.Println(err)
			}
		}

		if err := iter.Close(); err != nil {
			log.Fatal(err)
		}

		time.Sleep(10000 * time.Millisecond)
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

// go binary encoder
func ToGOB64(instrumentPrice common.InstrumentPrice) string {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(instrumentPrice)
	if err != nil {
		fmt.Println(`failed gob Encode`, err)
	}
	return base64.StdEncoding.EncodeToString(b.Bytes())
}
