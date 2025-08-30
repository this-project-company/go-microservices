package testing

import (
    "context"
    "testing"

    "go-microservices/customer-service/proto"
    "go-microservices/customer-service/controllers"
    "go-microservices/pkg/cache"
    "go-microservices/pkg/config"
)

func init() {
	cache.InitRedis()
    config.DBconnection()
}

func TestCustomer(t *testing.T) {
    s := &controllers.CustomerServer{} 

    req := &proto.GetCustomerRequest{Id: "1"}

    res, err := s.GetCustomer(context.Background(), req)
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    if res.Customer.Id != "1" {
        t.Errorf("expected id=1, got %s", res.Customer.Id)
    }

    if res.Customer.Name != "newman" {
        t.Errorf("expected name=newman, got %s", res.Customer.Name)
    }

    if res.Customer.Email != "newman@hotmail.com" {
        t.Errorf("expected email=newman@hotmail.com, got %s", res.Customer.Email)
    }
}
