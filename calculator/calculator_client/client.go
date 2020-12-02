package main

import (
	"../calculatorpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewSumServiceClient(cc)

	callSum(c)
}

func callSum(c calculatorpb.SumServiceClient) {
	req := &calculatorpb.SumRequest{
		FirstNum:  3,
		SecondNum: 10,
	}
	resp, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Sum: %v", err)
	}

	fmt.Println("Our sum is:", resp.Result)
}
