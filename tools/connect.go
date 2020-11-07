package main

import (
	"context"
	"flag"
	pb "github.com/billcchung/example-service/proto"
	"google.golang.org/grpc"
	"log"
)

func main() {
	var address string
	flag.StringVar(&address, "ADDRESS", ":8080", "The address to listen for connections")
	flag.Parse()
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect to address '%s', err: %s", address, err)
	}
	ctx := context.Background()
	client := pb.NewPingClient(conn)
	res, err := client.Get(ctx, &pb.PingRequest{Message_ID: "test", MessageBody: "test body"})
	log.Printf("Received response from server: %s \n", res)

	res, err = client.GetAfter(ctx, &pb.PingRequestWithSleep{Message_ID: "testSleep", MessageBody: "test body sleep", Sleep: 2})
	log.Printf("Received response from server: %s \n", res)

	res, err = client.GetRandom(ctx, &pb.PingRequest{Message_ID: "testRandom", MessageBody: "test body random"})
	log.Printf("Received response from server: %s \n", res)
}
