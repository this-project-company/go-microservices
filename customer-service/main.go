package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"go-microservices/customer-service/controllers"
	pb "go-microservices/customer-service/proto"
	"go-microservices/pkg/cache"
	"go-microservices/pkg/config"
	initializers "go-microservices/pkg/initializer"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Panic recovery interceptor
func recoveryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered in %s: %v", info.FullMethod, r)
			err = status.Errorf(codes.Internal, "internal server error")
		}
	}()
	return handler(ctx, req)
}


func init() {
	initializers.LoadEnvvariables()
}


func main() {
	// Start gRPC server with recovery interceptor
    cache.InitRedis()

	config.DBconnection()
	defer config.CloseDB()

	server := grpc.NewServer(
		grpc.UnaryInterceptor(recoveryInterceptor),
	)

	port := os.Getenv("CUSTOMER_PORT")

	pb.RegisterCustomerServiceServer(server, &controllers.CustomerServer{})

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Printf("🚀 Customer service running on %s\n", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
