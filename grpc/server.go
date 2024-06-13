package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	//set up connection and error handling
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("error with listening on port 9000: %v", err)
	}

	//import grpc server
	grpcServer := grpc.NewServer()

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("fail to serve gRPC server over port 9000 %v", err)
	}
}
