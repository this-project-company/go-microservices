package routes

import (
    "context"
    "log"
    "time"

    "github.com/gin-gonic/gin"
    pb "go-microservices/customer-service/proto"

    "google.golang.org/grpc"
)

func CustomerRoutes(rg *gin.RouterGroup) {
    customer := rg.Group("/customer")
    {
        customer.PUT("/edit/:id", editCustomer)
        customer.GET("/:id", getCustomer)
        customer.DELETE("/:id", deleteCustomer)
    }
}

func editCustomer(c *gin.Context) {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("could not connect: %v", err)
    }
    defer conn.Close()
    client := pb.NewCustomerServiceClient(conn)

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
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("could not connect: %v", err)
    }
    defer conn.Close()
    client := pb.NewCustomerServiceClient(conn)

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
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("could not connect: %v", err)
    }
    defer conn.Close()
    client := pb.NewCustomerServiceClient(conn)
    
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

