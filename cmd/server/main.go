package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"

	pb "github.com/kritika0598/simple-grpc/proto"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *server) SayHelloStream(req *pb.HelloStreamRequest, stream pb.Greeter_SayHelloStreamServer) error {
	log.Printf("Say hello %v times", req.GetTimes())
	//var wg sync.WaitGroup
	for i := 0; i < int(req.GetTimes()); i += 1 {
		time.Sleep(1 * time.Second)
		resp := pb.HelloReply{Message: fmt.Sprintf("Hello %s for %v time", req.GetName(), i)}
		if err := stream.Send(&resp); err != nil {
			log.Printf("send error %v", err)
		}
		log.Printf("finishing request for %d", i)
		//wg.Add(1)
		//go func(count int32) {
		//	defer wg.Done()
		//
		//}(int32(i))
	}
	//wg.Wait()
	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
