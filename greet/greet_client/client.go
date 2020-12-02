package main

import (
	"../greetpb"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

func main() {
	fmt.Println("hello im a client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	fmt.Println("created client", c)
}
