package main

import (
	"github.com/gocql/gocql"
	"github.com/gomodule/redigo/redis"
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
	// new Redis connection pool
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

			_, err := conn.Do("SET", instrumentKey, price)
			if err != nil {
				log.Println(err)
			}

			log.Println("Instrument:", instrumentKey, "Price:", price)
		}

		if err := iter.Close(); err != nil {
			log.Fatal(err)
		}

		time.Sleep(60000 * time.Millisecond)
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
