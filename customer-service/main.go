package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"go-microservices/customer-service/controllers"
	pb "go-microservices/customer-service/proto"
	"go-microservices/pkg/cache"
	"go-microservices/pkg/config"

	"github.com/joho/godotenv"
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

func LoadEnvvariables() {
    err := godotenv.Load()
    if err != nil {
        log.Println("⚠️ No .env file found, relying on system environment variables")
    }
}


func main() {
	// Start gRPC server with recovery interceptor
	LoadEnvvariables()
    cache.InitRedis()

	config.DBconnection()
	defer config.CloseDB()

	server := grpc.NewServer(
		grpc.UnaryInterceptor(recoveryInterceptor),
	)

	pb.RegisterCustomerServiceServer(server, &controllers.CustomerServer{})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("🚀 Customer service running on :50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
