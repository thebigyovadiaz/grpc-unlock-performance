package main

import (
	"context"
	"fmt"
	pb "github.com/thebigyovadiaz/grpc-unlock-performance/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

func main() {
	cred, err := credentials.NewClientTLSFromFile("server.crt", "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(cred))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "Peter", Lastname: "Crouch"})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			// This is a gRPC error
			fmt.Println("gRPC Error Code", st.Code())
			fmt.Println("gRPC Error Message", st.Message())
		} else {
			// Non gRPC error
			log.Fatalf("Non gRPC error: %v", err)
		}
	}

	log.Printf("Greeting: %s", r.GetMessage())
}
