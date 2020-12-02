package main

import (
	"../calculatorpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

func (s server) PrimeNumberDecomposition(request *calculatorpb.PrimeNumberRequest, decompositionServer calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	panic("implement me")
}

func (s server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Println("Sum Service Invoked")
	fn := req.GetFirstNum()
	sn := req.GetSecondNum()
	result := fn + sn
	return &calculatorpb.SumResponse{Result: result}, nil
}

func main() {
	fmt.Println("Starting up Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
