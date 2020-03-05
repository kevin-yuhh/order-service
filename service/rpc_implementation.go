package service

import (
	"time"

	"github.com/TRON-US/soter-order-service/charge"
	"github.com/TRON-US/soter-order-service/common/errorm"
	"github.com/TRON-US/soter-order-service/config"
	"github.com/TRON-US/soter-order-service/model"
	orderPb "github.com/TRON-US/soter-proto/order-service"

	"github.com/prometheus/client_golang/prometheus"
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
	t := time.Now()
	rpcRequestCount.With(prometheus.Labels{"method": "QueryBalanceTotal", "date": t.Format("2006-01-02")}).Inc()

	// Check input params.
	address := in.GetAddress()
	if address == "" {
		rpcRequestCount.With(prometheus.Labels{"method": "QueryBalanceFailed", "date": t.Format("2006-01-02")}).Inc()
		return nil, errorm.RequestParamEmpty
	}

	defer func(t time.Time) {
		rpcRequestDuration.With(prometheus.Labels{"method": "QueryBalance", "date": t.Format("2006-01-02")}).Observe(float64(time.Since(t).Microseconds()) / 1000)
	}(t)

	// Query ledger info by address.
	ledger, err := s.QueryBalanceController(address)
	if err != nil {
		rpcRequestCount.With(prometheus.Labels{"method": "QueryBalanceError", "date": t.Format("2006-01-02")}).Inc()
		return nil, err
	}

	rpcRequestCount.With(prometheus.Labels{"method": "QueryBalanceSuccess", "date": t.Format("2006-01-02")}).Inc()
	return &orderPb.QueryBalanceResponse{Balance: ledger.Balance}, nil
}

// Query order info by request id.
func (s *Server) QueryOrder(ctx context.Context, in *orderPb.QueryOrderRequest) (*orderPb.QueryOrderResponse, error) {
	t := time.Now()
	rpcRequestCount.With(prometheus.Labels{"method": "QueryOrderTotal", "date": t.Format("2006-01-02")}).Inc()

	// Check input params.
	address := in.GetAddress()
	requestId := in.GetRequestId()
	if address == "" || requestId == "" {
		rpcRequestCount.With(prometheus.Labels{"method": "QueryOrderFailed", "date": t.Format("2006-01-02")}).Inc()
		return nil, errorm.RequestParamEmpty
	}

	defer func(t time.Time) {
		rpcRequestDuration.With(prometheus.Labels{"method": "QueryOrder", "date": t.Format("2006-01-02")}).Observe(float64(time.Since(t).Microseconds()) / 1000)
	}(t)

	// Query order info by request id and address.
	order, err := s.QueryOrderController(requestId, address)
	if err != nil {
		rpcRequestCount.With(prometheus.Labels{"method": "QueryOrderError", "date": t.Format("2006-01-02")}).Inc()
		return nil, err
	}

	rpcRequestCount.With(prometheus.Labels{"method": "QueryOrderSuccess", "date": t.Format("2006-01-02")}).Inc()
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
	t := time.Now()
	rpcRequestCount.With(prometheus.Labels{"method": "CreateOrderTotal", "date": t.Format("2006-01-02")}).Inc()

	// Check input params.
	address := in.GetAddress()
	requestId := in.GetRequestId()
	fileName := in.GetFileName()
	fileSize := in.GetFileSize()
	if address == "" || requestId == "" || fileName == "" || fileSize <= 0 {
		rpcRequestCount.With(prometheus.Labels{"method": "CreateOrderFailed", "date": t.Format("2006-01-02")}).Inc()
		return nil, errorm.RequestParamEmpty
	}

	defer func(t time.Time) {
		defer rpcRequestDuration.With(prometheus.Labels{"method": "CreateOrder", "date": t.Format("2006-01-02")}).Observe(float64(time.Since(t).Microseconds()) / 1000)
	}(t)

	// Create order by address, requestId, fileName and fileSize.
	id, err := s.CreateOrderController(address, requestId, fileName, fileSize)
	if err != nil {
		rpcRequestCount.With(prometheus.Labels{"method": "CreateOrderError", "date": t.Format("2006-01-02")}).Inc()
		return nil, err
	}

	rpcRequestCount.With(prometheus.Labels{"method": "CreateOrderSuccess", "date": t.Format("2006-01-02")}).Inc()
	return &orderPb.CreateOrderResponse{OrderId: *id, SaveDays: int64(s.Time)}, nil
}

// Update file hash and session id.
func (s *Server) UpdateOrder(ctx context.Context, in *orderPb.UpdateOrderRequest) (*orderPb.Null, error) {
	t := time.Now()
	rpcRequestCount.With(prometheus.Labels{"method": "UpdateOrderTotal", "date": t.Format("2006-01-02")}).Inc()

	// Check input params.
	fileHash := in.GetFileHash()
	sessionId := in.GetSessionId()
	nodeIp := in.GetNodeIp()
	orderId := in.GetOrderId()
	if fileHash == "" || sessionId == "" || nodeIp == "" || orderId == 0 {
		rpcRequestCount.With(prometheus.Labels{"method": "UpdateOrderFailed", "date": t.Format("2006-01-02")}).Inc()
		return nil, errorm.RequestParamEmpty
	}

	defer func(t time.Time) {
		defer rpcRequestDuration.With(prometheus.Labels{"method": "UpdateOrder", "date": t.Format("2006-01-02")}).Observe(float64(time.Since(t).Microseconds()) / 1000)
	}(t)

	// Update order by fileHash, sessionId, nodeIp and orderId.
	err := s.UpdateOrderController(fileHash, sessionId, nodeIp, orderId)
	if err != nil {
		rpcRequestCount.With(prometheus.Labels{"method": "UpdateOrderError", "date": t.Format("2006-01-02")}).Inc()
		return nil, err
	}

	rpcRequestCount.With(prometheus.Labels{"method": "UpdateOrderSuccess", "date": t.Format("2006-01-02")}).Inc()
	return &orderPb.Null{}, nil
}

// Submit order by order Id.
func (s *Server) SubmitOrder(ctx context.Context, in *orderPb.SubmitOrderRequest) (*orderPb.Null, error) {
	t := time.Now()
	rpcRequestCount.With(prometheus.Labels{"method": "SubmitOrderTotal", "date": t.Format("2006-01-02")}).Inc()

	// Check input params.
	orderId := in.GetOrderId()
	fileHash := in.GetFileHash()
	result := "error"
	if orderId <= 0 {
		rpcRequestCount.With(prometheus.Labels{"method": "SubmitOrderFailed", "date": t.Format("2006-01-02")}).Inc()
		return nil, errorm.RequestParamEmpty
	}

	defer func(t time.Time) {
		defer rpcRequestDuration.With(prometheus.Labels{"method": "SubmitOrder", "date": t.Format("2006-01-02")}).Observe(float64(time.Since(t).Microseconds()) / 1000)
	}(t)

	// Submit order by orderId.
	err := s.SubmitOrderController(fileHash, result, orderId)
	if err != nil {
		rpcRequestCount.With(prometheus.Labels{"method": "SubmitOrderError", "date": t.Format("2006-01-02")}).Inc()
		return nil, err
	}

	rpcRequestCount.With(prometheus.Labels{"method": "SubmitOrderSuccess", "date": t.Format("2006-01-02")}).Inc()
	return &orderPb.Null{}, nil
}

// Close order by order id.
func (s *Server) CloseOrder(ctx context.Context, in *orderPb.CloseOrderRequest) (*orderPb.Null, error) {
	t := time.Now()
	rpcRequestCount.With(prometheus.Labels{"method": "CloseOrderTotal", "date": t.Format("2006-01-02")}).Inc()

	// Check input params.
	orderId := in.GetOrderId()
	if orderId <= 0 {
		rpcRequestCount.With(prometheus.Labels{"method": "CloseOrderFailed", "date": t.Format("2006-01-02")}).Inc()
		return nil, errorm.RequestParamEmpty
	}

	defer func(t time.Time) {
		defer rpcRequestDuration.With(prometheus.Labels{"method": "CloseOrder", "date": t.Format("2006-01-02")}).Observe(float64(time.Since(t).Microseconds()) / 1000)
	}(t)

	// Close order by order id.
	err := s.CloseOrderController(orderId)
	if err != nil {
		rpcRequestCount.With(prometheus.Labels{"method": "CloseOrderError", "date": t.Format("2006-01-02")}).Inc()
		return nil, err
	}

	rpcRequestCount.With(prometheus.Labels{"method": "CloseOrderSuccess", "date": t.Format("2006-01-02")}).Inc()
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
