package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/ndjordjevic/go-sb/internal/kafka_common"
	"log"
	"net/http"
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

	log.Println(instruments)

	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": instruments})
}

func fetchAllUsers(c *gin.Context) {
	var users []kafka_common.User
	m := map[string]interface{}{}

	iter := session.Query("SELECT company, email, first_name, last_name, password, address, city, country, accounts FROM users").Iter()
	for iter.MapScan(m) {
		user := kafka_common.User{
			Company:   m["company"].(string),
			Email:     m["email"].(string),
			FirstName: m["first_name"].(string),
			LastName:  m["last_name"].(string),
			Password:  m["password"].(string),
			Address:   m["address"].(string),
			City:      m["city"].(string),
			Country:   m["country"].(string),
		}

		for _, v := range m["accounts"].([]map[string]interface{}) {
			account := kafka_common.Account{
				Balance:  v["balance"].(float64),
				Currency: v["currency"].(string),
			}
			user.Accounts = append(user.Accounts, account)
		}
		users = append(users, user)

		m = map[string]interface{}{}
	}

	log.Println(users)

	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": users})
}

func createOrder(c *gin.Context) {

}
