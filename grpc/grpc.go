package grpc

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func RegisterGrpc(port *int, c func(s *grpc.Server)) {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	c(s)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
