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

func (s server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	n := req.GetNumber()
	primeNumberDecomposition(n, stream)
	return nil
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

func primeNumberDecomposition(n int32, steam calculatorpb.CalculatorService_PrimeNumberDecompositionServer) {
	k := int32(2)
	for n > 1 {
		if n % k == 0 {
			err := steam.Send(&calculatorpb.PrimeNumberResponse{Result: k}); if err != nil {
				log.Fatalf("Error while attempting to send stream: %v", err)
			}
			n = n/k
		} else {
			k++
		}
	}
}
