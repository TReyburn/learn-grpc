package main

import (
	"../greetpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct{}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fname := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		res := "Hello "+fname+" number "+ strconv.Itoa(i)
		resp := greetpb.GreetManyTimesResponse{Result: res}
		_ = stream.Send(&resp)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet call invoked: %v\n", req)
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
