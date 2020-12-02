package main

import (
	"../greetpb"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	fmt.Println("hello im a client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	// doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	doBiDirectionalStream(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Travis",
			LastName:  "Reyburn",
		} }
	resp, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet RPC: %v", resp.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Solveig",
			LastName:  "Delabroye",
		},
	}

	stream, err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling Server Streaming RPC: %v", err)
	}
	for {
		msg, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatalf("Error when streaming: %v", err)
		}
		log.Println("Response from Stream:", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	reqs := []*greetpb.LongGreetRequest{
		{Greeting: &greetpb.Greeting{
			FirstName: "Travis",
			LastName:  "Reyburn",
		}},
		{Greeting: &greetpb.Greeting{
			FirstName: "Travioli",
			LastName:  "Reyburini",
		}},
		{Greeting: &greetpb.Greeting{
			FirstName: "Solveig",
			LastName:  "Delabroye",
		}},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {log.Fatalf("Error was calling LongGreet: %v", err)}

	for _, req := range reqs {
		fmt.Println("Sending request", req)
		err := stream.Send(req)
		if err != nil {log.Fatalf("Error while streaming: %v", err)}
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {log.Fatalf("Error when receiving response: %v", err)}

	fmt.Println("Response:", res.GetResult())
}

func doBiDirectionalStream(c greetpb.GreetServiceClient) {
	reqs := []*greetpb.GreetEveryoneRequest{
		{Greeting: &greetpb.Greeting{
			FirstName: "Travis",
			LastName:  "Reyburn",
		}},
		{Greeting: &greetpb.Greeting{
			FirstName: "Travioli",
			LastName:  "Reyburini",
		}},
		{Greeting: &greetpb.Greeting{
			FirstName: "Solveig",
			LastName:  "Delabroye",
		}},
	}

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {log.Fatalf("Error was calling LongGreet: %v", err)}

	waitc := make(chan struct{})

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending message:", req)
			err := stream.Send(req)
			if err != nil {log.Fatalf("Error when streaming: %v", err)}
			time.Sleep(time.Second)
		}
		err := stream.CloseSend()
		if err != nil {log.Fatalf("Error when closing stream: %v", err)}
	} ()

	go func() {
		for {
			resp, err := stream.Recv()
			if errors.Is(err, io.EOF) {close(waitc); break}
			if err != nil {log.Fatalf("Error receiving stream data: %v", err)}
			fmt.Println("Received:", resp.GetResult())
		}
	} ()
	<-waitc
}
