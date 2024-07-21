package main

import (
	"context"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"

	pb "github.com/kritika0598/simple-grpc/proto"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx := context.Background()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())

	streamReq := &pb.HelloStreamRequest{Name: "Kritika", Times: 15}
	stream, err := c.SayHelloStream(ctx, streamReq)
	if err != nil {
		log.Fatalf("open stream error %v", err)
	}

	done := make(chan bool)
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			log.Printf("resp received: %s", resp.Message)
		}
	}()

	<-done
	log.Printf("finished")
}
