package main

import (
	"../greetpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fn := req.GetGreeting().GetFirstName()
	result := "Hello " + fn
	res := greetpb.GreetResponse{Result: result}
	return &res, nil
}

func main() {
	fmt.Println("hello im a server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
