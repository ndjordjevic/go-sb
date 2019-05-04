package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/golang/protobuf/ptypes/empty"
	orderpb "github.com/ndjordjevic/go-sb/api/order"
	pricepb "github.com/ndjordjevic/go-sb/api/price"
	"github.com/ndjordjevic/go-sb/internal/common"
	"google.golang.org/grpc"
	"gopkg.in/olivere/elastic.v7"
	"log"
	"net/http"
	"reflect"
	"time"
)

var session *gocql.Session
var orderHandlerServiceClient orderpb.OrderHandlerServiceClient
var priceServiceClient pricepb.PriceServiceClient

func init() {
	// connect to Cassandra cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "go_sb"
	session, _ = cluster.CreateSession()
	log.Println("Connected to Cassandra.")
}

func main() {
	defer session.Close()

	// grpc client connection
	orderClientConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := orderClientConn.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// grpc client connection
	priceClientConn, err := grpc.Dial("localhost:50071", grpc.WithInsecure())

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := priceClientConn.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// grpc service client
	orderHandlerServiceClient = orderpb.NewOrderHandlerServiceClient(orderClientConn)
	priceServiceClient = pricepb.NewPriceServiceClient(priceClientConn)

	// gin gonic routes
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

	router.GET("/api/v1/go-sb/orders/search", searchOrders)

	router.GET("/api/v1/go-sb/prices/", getPrices)

	_ = router.Run(":8010")
}

func fetchAllInstruments(c *gin.Context) {

	var instruments []common.Instrument
	m := map[string]interface{}{}

	iter := session.Query("SELECT instrument_key, market, isin, currency, short_name, long_name, status, expiration_date FROM instruments").Iter()
	for iter.MapScan(m) {
		instruments = append(instruments, common.Instrument{
			InstrumentKey:  m["instrument_key"].(string),
			Market:         m["market"].(string),
			ISIN:           m["isin"].(string),
			Currency:       m["currency"].(string),
			ShortName:      m["short_name"].(string),
			LongName:       m["long_name"].(string),
			ExpirationDate: m["expiration_date"].(time.Time),
			Status:         m["status"].(string),
			Price:          float32(0),
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

	req := &orderpb.HandleOrderRequest{
		Order: &orderpb.Order{
			Email:         order.Email,
			InstrumentKey: order.InstrumentKey,
			Currency:      order.Currency,
			Size:          order.Size,
			Price:         order.Price,
		},
		Action: orderpb.HandleOrderRequest_NEW,
	}

	res, err := orderHandlerServiceClient.HandleOrder(context.Background(), req)

	if err != nil {
		log.Fatal(err)
	}

	var message string
	var statusCode int

	if res.Response == orderpb.HandleOrderResponse_OK {
		message = "Order is created successfully"
		statusCode = http.StatusCreated
	} else {
		message = res.ErrorMessage
		statusCode = http.StatusNotAcceptable
	}

	c.JSON(statusCode, gin.H{"status": statusCode, "message": message})
}

func searchOrders(c *gin.Context) {
	searchTerm, _ := c.GetQuery("q")

	ctx := context.Background()

	client, err := elastic.NewClient()
	if err != nil {
		panic(err)
	}

	termQuery := elastic.NewQueryStringQuery(searchTerm)
	searchResult, err := client.Search().
		Index("order").
		Query(termQuery).
		Sort("Created", false).
		//From(0).Size(10).
		Pretty(true).
		Do(ctx)
	if err != nil {
		panic(err)
	}

	//fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	var orders []common.Order

	var order common.Order
	for _, item := range searchResult.Each(reflect.TypeOf(order)) {
		if order, ok := item.(common.Order); ok {
			orders = append(orders, order)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": orders})
}

func getPrices(c *gin.Context) {
	// Get prices
	res, err := priceServiceClient.RequestPrices(context.Background(), &empty.Empty{})

	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": res.Prices})
}
