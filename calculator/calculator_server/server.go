package main

import (
	"../calculatorpb"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type server struct{}

func (s server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	nums := make([]int32, 0)

	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatalf("Error while streaming: %v", err)
		}
		nums = append(nums, req.GetNumber())
	}

	res := float32(0.0)

	for _, num := range nums {
		res += float32(num)
	}

	res /= float32(len(nums))

	err := stream.SendAndClose(&calculatorpb.AverageResponse{Result: res})
	if err != nil {
		log.Fatalf("Error when sending response: %v", err)
	}

	return nil
}

func (s server) PrimeNumberDecomposition(
	req *calculatorpb.PrimeNumberRequest,
	stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer,
) error {
	fmt.Println("Prime Number Decomposition Service Invoked")
	n := req.GetNumber()
	primeNumberDecomposition(n, stream)
	return nil
}

func (s server) Sum(_ context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
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
		if n%k == 0 {
			err := steam.Send(&calculatorpb.PrimeNumberResponse{Result: k})
			if err != nil {
				log.Fatalf("Error while attempting to send stream: %v", err)
			}
			n /= k
		} else {
			k++
		}
	}
}
