package main

import (
	"../calculatorpb"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	//callSum(c)
	callPrimeNumber(c)
}

func callSum(c calculatorpb.CalculatorServiceClient) {
	req := calculatorpb.SumRequest{
		FirstNum:  25,
		SecondNum: 13,
	}

	resp, err := c.Sum(context.Background(), &req)
	if err != nil {
		log.Fatalf("Error while calling Sum: %v", err)
	}

	fmt.Println("Our sum is:", resp.GetResult())
}

func callPrimeNumber(c calculatorpb.CalculatorServiceClient) {
	n := int32(13467)
	req := calculatorpb.PrimeNumberRequest{Number: n}
	rs := make([]int32, 0)

	fmt.Println("Number:", n)

	stream, err := c.PrimeNumberDecomposition(context.Background(), &req)
	if err != nil {
		log.Fatalf("Error while calling Prime Num: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatalf("Error while streaming: %v", err)
		}
		rs = append(rs, resp.GetResult())
		fmt.Println("Number decomposes to:", resp.GetResult())
	}
	if len(rs) == 1 {
		fmt.Println(n, "is prime")
	} else {
		fmt.Println(n, "is not prime")
	}
	defer fmt.Println(n, "is equal to the following numbers multiplied together", rs)
}
