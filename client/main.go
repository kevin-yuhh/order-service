package main

import (
	"context"
	"fmt"

	orderPb "github.com/TRON-US/soter-order-service/proto"

	"github.com/gofrs/uuid"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:6661", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = conn.Close()
	}()
	client := orderPb.NewOrderServiceClient(conn)

	balance := QueryUserBalance(client, "TUsf2groYouQ7RzMkGcJH3PnSxFcwJCvrh")
	fmt.Printf("Balance is %v\n", balance)

	orderId := CreateOrder(client, "TUsf2groYouQ7RzMkGcJH3PnSxFcwJCvrh", 1000000)

	balance = QueryUserBalance(client, "TUsf2groYouQ7RzMkGcJH3PnSxFcwJCvrh")
	fmt.Printf("After create order, balance is %v\n", balance)

	CloseOrder(client, orderId)

	balance = QueryUserBalance(client, "TUsf2groYouQ7RzMkGcJH3PnSxFcwJCvrh")
	fmt.Printf("After close order, balance is %v\n", balance)

	orderId = CreateOrder(client, "TUsf2groYouQ7RzMkGcJH3PnSxFcwJCvrh", 1000000)

	balance = QueryUserBalance(client, "TUsf2groYouQ7RzMkGcJH3PnSxFcwJCvrh")
	fmt.Printf("After create order, balance is %v\n", balance)

	SubmitOrder(client, orderId)

	balance = QueryUserBalance(client, "TUsf2groYouQ7RzMkGcJH3PnSxFcwJCvrh")
	fmt.Printf("After submit order, balance is %v\n", balance)
}

func QueryUserBalance(c orderPb.OrderServiceClient, address string) int64 {
	request := &orderPb.QueryBalanceRequest{
		Address: address,
	}

	response, err := c.QueryBalance(context.Background(), request)
	if err != nil {
		panic(err)
	}

	return response.GetBalance()
}

func CreateOrder(c orderPb.OrderServiceClient, address string, fileSize int64) int64 {
	requestId, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	request := &orderPb.CreateOrderRequest{
		Address:   address,
		RequestId: requestId.String(),
		FileSize:  fileSize,
		FileName:  "test.txt",
	}

	response, err := c.CreateOrder(context.Background(), request)
	if err != nil {
		panic(err)
	}

	return response.GetOrderId()
}

func SubmitOrder(c orderPb.OrderServiceClient, orderId int64) {
	fileHash, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	request := &orderPb.SubmitOrderRequest{
		OrderId:  orderId,
		FileHash: fileHash.String(),
	}

	_, err = c.SubmitOrder(context.Background(), request)
	if err != nil {
		panic(err)
	}
}

func CloseOrder(c orderPb.OrderServiceClient, orderId int64) {
	request := &orderPb.CloseOrderRequest{
		OrderId: orderId,
	}

	_, err := c.CloseOrder(context.Background(), request)
	if err != nil {
		panic(err)
	}
}
