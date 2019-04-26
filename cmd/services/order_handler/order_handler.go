package main

import (
	"context"
	"github.com/gocql/gocql"
	orderpb "github.com/ndjordjevic/go-sb/api"
	"github.com/ndjordjevic/go-sb/internal/common"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type server struct{}

var session *gocql.Session

func init() {
	// connect to Cassandra cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "go_sb"
	session, _ = cluster.CreateSession()
	log.Println("Connected to Cassandra.")
}

func (*server) HandleOrder(ctx context.Context, req *orderpb.HandleOrderRequest) (*orderpb.HandleOrderResponse, error) {
	log.Println("New request to handle", req.Order)

	order := common.Order{
		Email:         req.Order.Email,
		InstrumentKey: req.Order.InstrumentKey,
		Currency:      req.Order.Currency,
		Size:          req.Order.Size,
		Price:         req.Order.Price,
	}

	resValidateOrder := callOrderValidatorService(order)

	res := &orderpb.HandleOrderResponse{}

	if resValidateOrder.Valid {
		log.Println("Order's valid")
		res.Response = orderpb.HandleOrderResponse_OK

		writeOrderToDBAsync(order)
	} else {
		log.Println(resValidateOrder.GetErrorMessage())
		res.Response = orderpb.HandleOrderResponse_ERROR
		res.ErrorMessage = resValidateOrder.GetErrorMessage()
	}

	return res, nil
}

func callOrderValidatorService(order common.Order) *orderpb.ValidateOrderResponse {
	reqValidateOrder := &orderpb.ValidateOrderRequest{
		Order: &orderpb.Order{
			Email:         order.Email,
			InstrumentKey: order.InstrumentKey,
			Currency:      order.Currency,
			Size:          order.Size,
			Price:         order.Price,
		},
	}

	// validation grpc client connection
	clientConn, err := grpc.Dial("localhost:50061", grpc.WithInsecure())

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := clientConn.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// grpc service client
	validateOrderServiceClient := orderpb.NewOrderValidatorServiceClient(clientConn)
	resValidateOrder, err := validateOrderServiceClient.ValidateOrder(context.Background(), reqValidateOrder)
	if err != nil {
		log.Fatal(err)
	}

	return resValidateOrder
}

func writeOrderToDBAsync(order common.Order) {
	go func() {
		order.Created = time.Now()
		order.UUID = gocql.TimeUUID()
		order.Status = "ACTIVE"

		// write order to Cassandra
		if err := session.Query(`INSERT INTO orders (uuid, email, instrument_key, currency, size, price, status, created) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			order.UUID, order.Email, order.InstrumentKey, order.Currency, order.Size, order.Price, order.Status, order.Created).Exec(); err != nil {
			log.Fatal(err)
		}
	}()
}

func main() {
	// order handler grpc service server
	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	orderpb.RegisterOrderHandlerServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
