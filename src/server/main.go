package main

import (
	"context"
	"fmt"
	pb "github.com/thebigyovadiaz/grpc-unlock-performance/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	if in.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "Name cannot be empty")
	}

	if in.Lastname == "" {
		return nil, status.Error(codes.InvalidArgument, "Lastname cannot be empty")
	}

	fullName := fmt.Sprintf("%s %s", in.GetName(), in.GetLastname())
	return &pb.HelloResponse{Message: fullName}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	cred, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
	if err != nil {
		log.Fatalf("Failed to generate credentials %v", err)
	}

	s := grpc.NewServer(grpc.Creds(cred))
	pb.RegisterGreeterServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
