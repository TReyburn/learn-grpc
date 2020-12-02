package main

import (
	"../greetpb"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct{}

func (s *server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("Receieved a streaming request on GreetEveryone")

	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {return nil}
		if err != nil {log.Fatalf("Error while reading client stream: %v", err)}
		fname := req.GetGreeting().GetFirstName()
		res := "Hello "+fname+"!"
		err = stream.Send(&greetpb.GreetEveryoneResponse{Result: res})
		if err != nil {log.Fatalf("Error when streaming to client: %v", err)}
	}
}

func (s *server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("Received a streaming response on LongGreet")
	res := ""
	for {
		msg, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			err := stream.SendAndClose(&greetpb.LongGreetResponse{Result: res})
			if err != nil {log.Fatalf("Error sending response: %v", err)}
			break
		}
		if err != nil {
			log.Fatalf("Error while streaming: %v", err)
		}
		fname := msg.GetGreeting().GetFirstName()
		res += "Hello "+fname+"! "
	}
	return nil
}

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

func (s *server) Greet(_ context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
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
