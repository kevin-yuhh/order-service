package main

import (
	"context"
	"fmt"

	orderPb "github.com/TRON-US/soter-order-service/proto"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:6661", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer cc.Close()
	c := orderPb.NewOrderServiceClient(cc)

	QueryUserBalance(c)
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
