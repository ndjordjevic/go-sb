package main

import (
	"context"
	orderpb "github.com/ndjordjevic/go-sb/api"
	"github.com/ndjordjevic/go-sb/internal/common"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

func (*server) HandleOrder(ctx context.Context, req *orderpb.HandleOrderRequest) (*orderpb.HandleOrderResponse, error) {
	order := common.Order{
		Email:         req.Order.Email,
		InstrumentKey: req.Order.InstrumentKey,
		Currency:      req.Order.Currency,
		Size:          req.Order.Size,
		Price:         req.Order.Price,
	}

	log.Println(order)

	res := &orderpb.HandleOrderResponse{
		Response: orderpb.HandleOrderResponse_OK,
	}

	return res, nil
}

func main() {
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
