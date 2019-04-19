package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/ndjordjevic/go-sb/internal/kafka_common"
	"log"
	"time"
)

var session *gocql.Session

func init() {
	// connect to Cassandra cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "go_sb"
	session, _ = cluster.CreateSession()
	log.Println("Connected to Cassandra.")
}

func main() {
	defer session.Close()

	router := gin.Default()

	instrumentsV1 := router.Group("/api/v1/go-sb/instruments/")
	instrumentsV1.GET("/", fetchAllInstruments)

	usersV1 := router.Group("/api/v1/go-sb/users/")
	usersV1.GET("/", fetchAllUsers)

	ordersV1 := router.Group("/api/v1/go-sb/orders/")
	ordersV1.POST("/", createOrder)

	_ = router.Run()
}

func fetchAllInstruments(c *gin.Context) {

	var instruments []kafka_common.Instrument
	m := map[string]interface{}{}

	iter := session.Query("SELECT market, isin, currency, short_name, long_name, status, expiration_date FROM instruments").Iter()
	for iter.MapScan(m) {
		instruments = append(instruments, kafka_common.Instrument{
			Market:         m["market"].(string),
			ISIN:           m["isin"].(string),
			Currency:       m["currency"].(string),
			ShortName:      m["short_name"].(string),
			LongName:       m["long_name"].(string),
			ExpirationDate: m["expiration_date"].(time.Time),
			Status:         m["status"].(string),
		})

		m = map[string]interface{}{}
	}

	fmt.Println(instruments)

	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
}

func fetchAllUsers(c *gin.Context) {

}

func createOrder(c *gin.Context) {

}
