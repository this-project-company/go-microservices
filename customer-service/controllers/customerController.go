package controllers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	pb "go-microservices/customer-service/proto"
	"go-microservices/pkg/cache"
	"go-microservices/pkg/config"
	"go-microservices/pkg/rabbitmq"

	"github.com/streadway/amqp"
)

// CustomerServer implements gRPC service
type CustomerServer struct {
	pb.UnimplementedCustomerServiceServer
}

// Exmaple : File sending like csv to bulk upload in DB (map the field corresponding to the DB cols)
func (s *CustomerServer) SendFile(ctx context.Context, req *pb.UploadCSVRequest) (*pb.UploadCSVResponse, error) {
    // Use a scanner to read the CSV line by line
    scanner := bufio.NewScanner(bytes.NewReader(req.FileContent))

    // store header as headers[]
    if !scanner.Scan() {
        return &pb.UploadCSVResponse{
            RowsProcessed: 0,
            Message:       "Empty CSV file",
        }, nil
    }
    headerLine := scanner.Text()
    headers := strings.Split(headerLine, ",")
    for i := range headers {
        headers[i] = strings.TrimSpace(headers[i]) // delete the header title
    }

    rowsProcessed := 0
    rowsSkipped := 0

    // Read remaining lines from the beginning 
    for scanner.Scan() {
        line := scanner.Text()

        // Skip completely empty lines
        if strings.TrimSpace(line) == "" {
            continue
        }

        fields := strings.Split(line, ",")
        rowData := make(map[string]string)
        for i, h := range headers {
            if i < len(fields) {
                rowData[h] = fields[i]
            } else {
                rowData[h] = "" // empty if missing
            }
        }

        rowsProcessed++
        fmt.Printf("Row %d: %+v\n", rowsProcessed, rowData) // debug log
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Error reading CSV:", err)
        rowsSkipped++ 
    }

    return &pb.UploadCSVResponse{
        RowsProcessed: int32(rowsProcessed),
        Message:       fmt.Sprintf("CSV processed successfully (skipped %d bad rows)", rowsSkipped),
    }, nil
}




// Example: GetCustomer (with Redis caching)
func (s *CustomerServer) GetCustomer(ctx context.Context, req *pb.GetCustomerRequest) (*pb.CustomerResponse, error) {
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

func (s * CustomerServer) CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.MessageOnlyResponse, error) {

    fmt.Println(req.Id, req.Name, req.Email)
    _, err := config.DBPool.Exec(
    ctx,
    `INSERT INTO customer (id, name, email) VALUES ($1, $2, $3)`, req.Id,
    req.Name, req.Email,)

    if err != nil {
        fmt.Println("error in insert")
        return nil, err
    }

    idStr := strconv.FormatInt(req.Id, 10)

    newUser := map[string]string{
        "id":   idStr,
        "name": req.Name,
        "email": req.Email,
    }
    body, _ := json.Marshal(newUser)

        err = rabbitmq.RabbitCh.Publish(
        "",             // exchange
        "user.created", // routing key
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )
    if err != nil {
        fmt.Println("error in insert")
        return nil, err    }

	return &pb.MessageOnlyResponse{
        Message:  "Customer Created Successfully",
    }, nil
}

//gorm migration