package main

import (
	"context"
	orderpb "github.com/ndjordjevic/go-sb/api"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
)

type validateOrderServer struct{}

func (*validateOrderServer) ValidateOrder(ctx context.Context, req *orderpb.ValidateOrderRequest) (*orderpb.ValidateOrderResponse, error) {
	log.Println("New request to validate", req.Order)

	res := &orderpb.ValidateOrderResponse{
		Valid: randomBool(),
	}

	if res.Valid == false {
		res.ErrorMessage = "Order didn't pass validation"
	}

	return res, nil

}

func randomBool() bool {
	return rand.Float32() < 0.5
}

func main() {
	// validate order grpc service server
	lis, err := net.Listen("tcp", "localhost:50061")

	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	orderpb.RegisterValidateOrderServiceServer(s, &validateOrderServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
