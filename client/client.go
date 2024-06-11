package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/serviceconfig"

	pb "projekt/proto"
)

var wg sync.WaitGroup

func loadRetryPolicy(file string) grpc.DialOption {
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("failed to read retry policy file: %v", err)
	}

	var cfg serviceconfig.ParseResult
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("failed to parse retry policy: %v", err)
	}

	return grpc.WithDefaultServiceConfig(string(data))
}

func clientAction(conn *grpc.ClientConn) {
	client := pb.NewSemaphoreClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	acquireResp, err := client.Acquire(ctx, &pb.AcquireRequest{Permits: 4})
	if err != nil {
		log.Fatalf("Could not acquire: %v", err)
	}

	if acquireResp.Success {
		log.Printf("Acquire success: %v", acquireResp.Success)
	}

	time.Sleep(100 * time.Millisecond)
	releaseResp, err := client.Release(ctx, &pb.ReleaseRequest{Permits: 4})
	if err != nil {
		log.Fatalf("could not release: %v", err)
	}
	log.Printf("Release success: %v", releaseResp.Success)
	wg.Done()
}

func main() {
	hostname := os.Args[1]
	port := os.Args[2]
	fullhost := fmt.Sprintf("%s:%s", hostname, port)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(fullhost, opts...)
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			clientAction(conn)
		}()
	}
	wg.Wait()
}
