package main

import (
	"context"
	"github.com/gocql/gocql"
	orderpb "github.com/ndjordjevic/go-sb/api/order"
	"google.golang.org/grpc"
	"log"
	"net"
)

type validateOrderServer struct{}

var session *gocql.Session

func init() {
	// connect to Cassandra cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "go_sb"
	session, _ = cluster.CreateSession()
	log.Println("Connected to Cassandra.")
}

func (*validateOrderServer) ValidateOrder(ctx context.Context, req *orderpb.ValidateOrderRequest) (*orderpb.ValidateOrderResponse, error) {
	log.Println("New request to validate", req.Order)

	res := &orderpb.ValidateOrderResponse{
		Valid: checkOrder(req),
	}

	if res.Valid == false {
		res.ErrorMessage = "Order didn't pass validation"
	}

	return res, nil

}

func checkOrder(req *orderpb.ValidateOrderRequest) bool {
	var accounts []map[string]interface{}
	if err := session.Query("SELECT accounts FROM users WHERE email = ? ALLOW FILTERING", req.Order.Email).Scan(&accounts); err != nil {
		log.Println(err)
	}
	var valid = false
	for _, v := range accounts {
		if v["currency"].(string) == req.Order.Currency {
			if float64(req.Order.Size*req.Order.Price) < v["balance"].(float64) {
				valid = true
			}
		}
	}
	return valid
}

func main() {
	// validate order grpc service server
	lis, err := net.Listen("tcp", "localhost:50061")

	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	orderpb.RegisterOrderValidatorServiceServer(s, &validateOrderServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
