package main

import (
	"log"
	"net"
	pb "projekt/proto"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp4", "0.0.0.0:8050")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSemaphoreServer(grpcServer, NewSemaphoreServer(10))
	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
