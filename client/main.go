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

	CreateOrder(c)
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

func CreateOrder(c orderPb.OrderServiceClient) {
	requestId, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	request := &orderPb.CreateOrderRequest{
		Address:   "TUsf2groYouQ7RzMkGcJH3PnSxFcwJCvrh",
		RequestId: requestId.String(),
		Amount:    100,
	}

	response, err := c.CreateOrder(context.Background(), request)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.GetOrderId())
}
