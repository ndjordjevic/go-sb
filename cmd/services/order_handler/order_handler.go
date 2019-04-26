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
	log.Println("New request to handle", req.Order)

	order := common.Order{
		Email:         req.Order.Email,
		InstrumentKey: req.Order.InstrumentKey,
		Currency:      req.Order.Currency,
		Size:          req.Order.Size,
		Price:         req.Order.Price,
	}

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
	validateOrderServiceClient := orderpb.NewValidateOrderServiceClient(clientConn)
	resValidateOrder, err := validateOrderServiceClient.ValidateOrder(context.Background(), reqValidateOrder)
	if err != nil {
		log.Fatal(err)
	}

	res := &orderpb.HandleOrderResponse{}

	if resValidateOrder.Valid {
		log.Println("Order's valid")
		res.Response = orderpb.HandleOrderResponse_OK

		// persist order ...
	} else {
		log.Println(resValidateOrder.GetErrorMessage())
		res.Response = orderpb.HandleOrderResponse_ERROR
		res.ErrorMessage = resValidateOrder.GetErrorMessage()
	}

	return res, nil
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
