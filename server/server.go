package main

import (
	"context"
	"log"
	"sync"

	pb "projekt/proto"
)

type semaphoreServer struct {
	pb.UnimplementedSemaphoreServer
	permits int32
	mutex   sync.Mutex
	cond    *sync.Cond
}

func NewSemaphoreServer(permits int32) *semaphoreServer {
	s := &semaphoreServer{
		permits: permits,
	}
	s.cond = sync.NewCond(&s.mutex)
	return s
}

func (s *semaphoreServer) Acquire(ctx context.Context, req *pb.AcquireRequest) (*pb.AcquireResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for s.permits < req.Permits {
		s.cond.Wait()
	}
	s.permits -= req.Permits
	log.Printf("Acquired %d permits. Current permits %d", req.Permits, s.permits)
	return &pb.AcquireResponse{Success: true}, nil
}

func (s *semaphoreServer) Release(ctx context.Context, req *pb.ReleaseRequest) (*pb.ReleaseResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.permits += req.Permits
	s.cond.Signal()
	log.Printf("Released %d permits. Current permits %d", req.Permits, s.permits)
	return &pb.ReleaseResponse{Success: true}, nil
}
