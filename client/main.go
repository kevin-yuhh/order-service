package main

import (
	"context"
	"fmt"

	orderPb "github.com/TRON-US/soter-order-service/proto"

	"github.com/gofrs/uuid"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:6661", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer cc.Close()
	c := orderPb.NewOrderServiceClient(cc)

	orderId := CreateOrder(c, "TUsf2groYouQ7RzMkGcJH3PnSxFcwJCvrh", 1000000)

	fmt.Println(orderId)

	SubmitOrder(c, orderId)
}

func QueryUserBalance(c orderPb.OrderServiceClient) {
	request := &orderPb.QueryBalanceRequest{
		Address: "TUsf2groYouQ7RzMkGcJH3PnSxFcwJCvrh",
	}

	response, err := c.QueryBalance(context.Background(), request)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.GetBalance())
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
