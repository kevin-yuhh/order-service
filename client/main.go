package main

import (
	"context"
	"fmt"

	orderPb "github.com/TRON-US/soter-proto/order-service"

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

	address := "TUsf2groYouQ7RzMkGcJH3PnSxFcwJCvrh"

	balance := QueryUserBalance(client, address)
	fmt.Printf("Balance is %v\n", balance)

	orderId, requestId := CreateOrder(client, address, 1000000)

	balance = QueryUserBalance(client, address)
	fmt.Printf("After create order, balance is %v\n", balance)

	fmt.Println(QueryOrderInfo(client, address, requestId))

	CloseOrder(client, orderId)

	fmt.Println(QueryOrderInfo(client, address, requestId))

	balance = QueryUserBalance(client, address)
	fmt.Printf("After close order, balance is %v\n", balance)

	orderId, requestId = CreateOrder(client, address, 1000000)

	balance = QueryUserBalance(client, address)
	fmt.Printf("After create order, balance is %v\n", balance)

	fmt.Println(QueryOrderInfo(client, address, requestId))

	SubmitOrder(client, orderId)

	fmt.Println(QueryOrderInfo(client, address, requestId))

	balance = QueryUserBalance(client, address)
	fmt.Printf("After submit order, balance is %v\n", balance)

	orderId, _ = PrepareRenew(client, 1)

	balance = QueryUserBalance(client, address)
	fmt.Printf("After prepare renew, balance is %v\n", balance)

	CloseOrder(client, orderId)

	balance = QueryUserBalance(client, address)
	fmt.Printf("After close renew order, balance is %v\n", balance)

	orderId, _ = PrepareRenew(client, 1)

	balance = QueryUserBalance(client, address)
	fmt.Printf("After prepare renew, balance is %v\n", balance)

	SubmitOrder(client, orderId)

	balance = QueryUserBalance(client, address)
	fmt.Printf("After submit renew order, balance is %v\n", balance)
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

func QueryOrderInfo(c orderPb.OrderServiceClient, address, requestId string) *orderPb.QueryOrderResponse {
	request := &orderPb.QueryOrderRequest{
		Address:   address,
		RequestId: requestId,
	}

	response, err := c.QueryOrder(context.Background(), request)
	if err != nil {
		panic(err)
	}

	return response
}

func CreateOrder(c orderPb.OrderServiceClient, address string, fileSize int64) (int64, string) {
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

	return response.GetOrderId(), requestId.String()
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

func PrepareRenew(c orderPb.OrderServiceClient, fileId int64) (int64, int64) {
	request := &orderPb.PrepareRenewRequest{
		FileId: fileId,
	}

	response, err := c.PrepareRenew(context.Background(), request)
	if err != nil {
		panic(err)
	}

	return response.GetOrderId(), response.GetStatus()
}
