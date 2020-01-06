package service

import (
	"github.com/TRON-US/soter-order-service/charge"
	"github.com/TRON-US/soter-order-service/common/errorm"
	"github.com/TRON-US/soter-order-service/config"
	"github.com/TRON-US/soter-order-service/model"
	orderPb "github.com/TRON-US/soter-proto/order-service"

	"golang.org/x/net/context"
)

type Server struct {
	DbConn *model.Database
	Fee    *charge.FeeCalculator
	Config *config.Configuration
	Time   int
}

// Query balance by address.
func (s *Server) QueryBalance(ctx context.Context, in *orderPb.QueryBalanceRequest) (*orderPb.QueryBalanceResponse, error) {
	// Check input params.
	address := in.GetAddress()
	if address == "" {
		return nil, errorm.RequestParamEmpty
	}

	// Query ledger info by address.
	ledger, err := s.QueryBalanceController(address)
	if err != nil {
		return nil, err
	}

	return &orderPb.QueryBalanceResponse{Balance: ledger.Balance}, nil
}

// Query order info by request id.
func (s *Server) QueryOrder(ctx context.Context, in *orderPb.QueryOrderRequest) (*orderPb.QueryOrderResponse, error) {
	// Check input params.
	address := in.GetAddress()
	requestId := in.GetRequestId()
	if address == "" || requestId == "" {
		return nil, errorm.RequestParamEmpty
	}

	// Query order info by request id and address.
	order, err := s.QueryOrderController(requestId, address)
	if err != nil {
		return nil, err
	}

	return &orderPb.QueryOrderResponse{
		Type:        order.OrderType,
		FileName:    order.FileName,
		FileSize:    order.FileSize,
		FileHash:    order.FileHash,
		Fee:         order.Amount,
		Status:      order.Status,
		Description: "",
	}, nil
}

// Create order by address and requestId.
func (s *Server) CreateOrder(ctx context.Context, in *orderPb.CreateOrderRequest) (*orderPb.CreateOrderResponse, error) {
	// Check input params.
	address := in.GetAddress()
	requestId := in.GetRequestId()
	fileName := in.GetFileName()
	fileSize := in.GetFileSize()
	if address == "" || requestId == "" || fileName == "" || fileSize <= 0 {
		return nil, errorm.RequestParamEmpty
	}

	// Create order by address, requestId, fileName and fileSize.
	id, err := s.CreateOrderController(address, requestId, fileName, fileSize)
	if err != nil {
		return nil, err
	}

	return &orderPb.CreateOrderResponse{OrderId: *id}, nil
}

// Update file hash and session id.
func (s *Server) UpdateOrder(ctx context.Context, in *orderPb.UpdateOrderRequest) (*orderPb.Null, error) {
	// Check input params.
	fileHash := in.GetFileHash()
	sessionId := in.GetSessionId()
	nodeIp := in.GetNodeIp()
	orderId := in.GetOrderId()
	if fileHash == "" || sessionId == "" || nodeIp == "" || orderId == 0 {
		return nil, errorm.RequestParamEmpty
	}

	// Update order by fileHash, sessionId, nodeIp and orderId.
	err := s.UpdateOrderController(fileHash, sessionId, nodeIp, orderId)
	if err != nil {
		return nil, err
	}

	return &orderPb.Null{}, nil
}

// Submit order by order Id.
func (s *Server) SubmitOrder(ctx context.Context, in *orderPb.SubmitOrderRequest) (*orderPb.Null, error) {
	// Check input params.
	orderId := in.GetOrderId()
	fileHash := in.GetFileHash()
	result := "error"
	if orderId <= 0 {
		return nil, errorm.RequestParamEmpty
	}

	// Submit order by orderId.
	err := s.SubmitOrderController(fileHash, result, orderId)
	if err != nil {
		return nil, err
	}

	return &orderPb.Null{}, nil
}

// Close order by order id.
func (s *Server) CloseOrder(ctx context.Context, in *orderPb.CloseOrderRequest) (*orderPb.Null, error) {
	// Check input params.
	orderId := in.GetOrderId()
	if orderId <= 0 {
		return nil, errorm.RequestParamEmpty
	}

	// Close order by order id.
	err := s.CloseOrderController(orderId)
	if err != nil {
		return nil, err
	}

	return &orderPb.Null{}, nil
}

// Prepare renew.
func (s *Server) PrepareRenew(ctx context.Context, in *orderPb.PrepareRenewRequest) (*orderPb.PrepareRenewResponse, error) {
	// Check input params.
	fileId := in.GetFileId()
	if fileId <= 0 {
		return nil, errorm.RequestParamEmpty
	}

	return &orderPb.PrepareRenewResponse{}, nil
}
