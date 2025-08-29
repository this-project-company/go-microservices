package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pb "go-microservices/customer-service/proto"
	"go-microservices/pkg/cache"
	"go-microservices/pkg/config"
)

// CustomerServer implements gRPC service
type CustomerServer struct {
	pb.UnimplementedCustomerServiceServer
}

func defaultIfEmpty(val, def string) string {
    if val == "" {
        return def
    }
    return val
}


// Example: GetCustomer (with Redis caching)
func (s *CustomerServer) GetCustomer(ctx context.Context, req *pb.GetCustomerRequest) (*pb.CustomerResponse, error) {
	fmt.Println("Entered Get")
    key := fmt.Sprintf("customer:%s", req.Id)

    // 1. Try Redis cache
    cached, err := cache.Get(key)
    if err == nil && cached != "" {
        var customer pb.Customer
        if jsonErr := json.Unmarshal([]byte(cached), &customer); jsonErr == nil {
            return &pb.CustomerResponse{
                Customer: &customer,
                Message:  "From Redis cache",
            }, nil
        }
    }                   

var name, email string
err = config.DBPool.QueryRow(
    ctx,
    `SELECT name, email FROM customer WHERE id=$1`,
    req.Id,
).Scan(&name, &email)

if err != nil {
    return nil, err
}

customer := &pb.Customer{
    Id:    req.Id,
    Name:  name,
    Email: email,
}

    // 3. Save to Redis
    data, _ := json.Marshal(customer)
    _ = cache.Set(key, string(data), 5*time.Minute)

    return &pb.CustomerResponse{
        Customer: customer,
        Message:  "From DB + cached in Redis",
    }, nil
}


func (s *CustomerServer) DeleteCustomer(ctx context.Context, req *pb.DeleteCustomerRequest) (*pb.MessageOnlyResponse, error) {
	fmt.Println("Entered delete")	
	key := fmt.Sprintf("customer:%s", req.Id)
	fmt.Println(key)

		cached, err := cache.Get(key)
		if err != nil && cached == "" {
			return &pb.MessageOnlyResponse{
				Message:  "No user",
				}, nil	
		}
		
	if err := cache.Delete(key); err != nil {
		return &pb.MessageOnlyResponse{
            Message:  "Customer Delete From Redis cache",
        }, nil
	}
		
		return &pb.MessageOnlyResponse{
            Message:  "Customer Delete From Redis cache",
        }, nil

}

//gorm migration