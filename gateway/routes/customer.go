package routes

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "go-microservices/customer-service/proto"

	"github.com/gin-gonic/gin"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CustomerRoutes(rg *gin.RouterGroup) {
    customer := rg.Group("/customer")
    {
        customer.PUT("/edit/:id", editCustomer)
        customer.GET("/:id", getCustomer)
        customer.DELETE("/:id", deleteCustomer)
    }
}

func newCustomerClient() pb.CustomerServiceClient {
    customerPort := os.Getenv("CUSTOMER_PORT")
    address := fmt.Sprintf("localhost%s", customerPort)

    conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("could not connect: %v", err)
    }
    // ⚠️ Don’t close conn here! Return client & keep conn open
    return pb.NewCustomerServiceClient(conn)
}


func editCustomer(c *gin.Context) {

    client := newCustomerClient()

    id := c.Param("id")
    var body struct {
        Name string `json:"name"`
    }
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(400, gin.H{"error": "invalid input"})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
    defer cancel()

    res, err := client.EditCustomer(ctx, &pb.EditCustomerRequest{
        Id:   id,
        Name: body.Name,
    })
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, res)
}

func getCustomer(c *gin.Context) {

    client := newCustomerClient()

    id := c.Param("id")

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
    defer cancel()

    res, err := client.GetCustomer(ctx, &pb.GetCustomerRequest{Id: id})
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, res)
}

func deleteCustomer(c *gin.Context) {

    client := newCustomerClient()
    
    id := c.Param("id")

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
    defer cancel()

    res, err := client.DeleteCustomer(ctx, &pb.DeleteCustomerRequest{Id: id})
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, res)
}

