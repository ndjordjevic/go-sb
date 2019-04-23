package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/ndjordjevic/go-sb/internal/common"
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

	// CORS for https://foo.com and https://github.com origins, allowing:
	// - PUT and PATCH methods
	// - Origin header
	// - Credentials share
	// - Preflight requests cached for 12 hours
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		//AllowOriginFunc: func(origin string) bool {
		//	return origin == "http://localhost:8010"
		//},
		MaxAge: 12 * time.Hour,
	}))

	instrumentsV1 := router.Group("/api/v1/go-sb/instruments/")
	instrumentsV1.GET("/", fetchAllInstruments)

	usersV1 := router.Group("/api/v1/go-sb/users/")
	usersV1.GET("/", fetchAllUsers)

	ordersV1 := router.Group("/api/v1/go-sb/orders/")
	ordersV1.POST("/", createOrder)

	_ = router.Run(":8010")
}

func fetchAllInstruments(c *gin.Context) {

	var instruments []common.Instrument
	m := map[string]interface{}{}

	iter := session.Query("SELECT market, isin, currency, short_name, long_name, status, expiration_date FROM instruments").Iter()
	for iter.MapScan(m) {
		instruments = append(instruments, common.Instrument{
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
	var users []common.User
	m := map[string]interface{}{}

	iter := session.Query("SELECT company, email, first_name, last_name, password, address, city, country, accounts FROM users").Iter()
	for iter.MapScan(m) {
		user := common.User{
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
			account := common.Account{
				Balance:  v["balance"].(float64),
				Currency: v["currency"].(string),
			}
			user.Accounts = append(user.Accounts, account)
		}
		users = append(users, user)

		m = map[string]interface{}{}
	}

	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": users})
}

func createOrder(c *gin.Context) {
	var order common.Order

	if err := c.BindJSON(&order); err != nil {
		log.Fatal(err)
	}

	order.Created = time.Now()
	order.UUID = gocql.TimeUUID()

	// Set to ACTIVE if it's valid
	order.Status = "ACTIVE"

	// write order to Cassandra
	if err := session.Query(`INSERT INTO orders (uuid, email, instrument_key, currency, size, price, status, created) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		order.UUID, order.Email, order.InstrumentKey, order.Currency, order.Size, order.Price, order.Status, order.Created).Exec(); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Order is created successfully!", "resourceId": order.UUID})
}
