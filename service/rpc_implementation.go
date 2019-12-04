package service

import (
	"time"

	"github.com/TRON-US/soter-order-service/charge"
	"github.com/TRON-US/soter-order-service/common/constants"
	"github.com/TRON-US/soter-order-service/common/errorm"
	"github.com/TRON-US/soter-order-service/config"
	"github.com/TRON-US/soter-order-service/model"
	orderPb "github.com/TRON-US/soter-order-service/proto"

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
	ledger, err := s.DbConn.QueryLedgerInfoByAddress(address)
	if err != nil {
		return nil, err
	}

	return &orderPb.QueryBalanceResponse{Balance: ledger.Balance}, nil
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

	// Query ledger info by address.
	ledger, err := s.DbConn.QueryLedgerInfoByAddress(address)
	if err != nil {
		return nil, err
	}

	// Calculate fee of this order.
	amount := s.Fee.Fee(fileSize, ledger.TotalTimes, s.Time)

	// Check balance illegal.
	if ledger.Balance < amount {
		return nil, errorm.InsufficientBalance
	}

	// Open transaction.
	session := s.DbConn.DB.NewSession()
	err = session.Begin()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	// Insert file information
	fileId, err := model.InsertFileInfo(session, ledger.UserId, fileSize, fileName, int(time.Now().Local().Unix())+s.Time*86400)
	if err != nil {
		_ = session.Rollback()
		return nil, err
	}

	// Insert order information.
	id, err := model.InsertOrderInfo(session, ledger.UserId, fileId, amount, s.Fee.StrategyId, requestId, s.Time)
	if err != nil {
		_ = session.Rollback()
		return nil, err
	}

	// Freeze user balance.
	err = model.UpdateUserBalance(session, ledger.Balance-amount, ledger.FreezeBalance+amount, ledger.Version, ledger.Id, ledger.Address, int(time.Now().Local().Unix()))
	if err != nil {
		_ = session.Rollback()
		return nil, err
	}

	// Submit transaction.
	err = session.Commit()
	if err != nil {
		return nil, err
	}

	return &orderPb.CreateOrderResponse{OrderId: id}, nil
}

// Submit order by order Id.
func (s *Server) SubmitOrder(ctx context.Context, in *orderPb.SubmitOrderRequest) (*orderPb.SubmitOrderResponse, error) {
	// Check input params.
	orderId := in.GetOrderId()
	fileHash := in.GetFileHash()
	if orderId <= 0 || fileHash == "" {
		return nil, errorm.RequestParamEmpty
	}

	// Get order info by order id.
	order, err := s.DbConn.QueryOrderInfoById(orderId)
	if err != nil {
		return nil, err
	}

	// Query ledger info by address.
	ledger, err := s.DbConn.QueryLedgerInfoByAddress(order.Address)
	if err != nil {
		return nil, err
	}

	// Check freeze balance illegal.
	if ledger.FreezeBalance < order.Amount {
		return nil, errorm.InsufficientBalance
	}

	// Open transaction.
	session := s.DbConn.DB.NewSession()
	err = session.Begin()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	// Update file hash by file id.
	err = model.UpdateFileHash(session, order.FileId, fileHash)
	if err != nil {
		_ = session.Rollback()
		return nil, err
	}

	// Update ledger information by ledger id.
	err = model.UpdateLedgerInfo(session, ledger.TotalSize+order.FileSize, ledger.Balance, ledger.FreezeBalance-order.Amount, ledger.TotalFee+order.Amount, ledger.Version, ledger.Id, order.Address, int(time.Now().Local().Unix()))
	if err != nil {
		_ = session.Rollback()
		return nil, err
	}

	// Update order status by order id.
	err = model.UpdateOrderStatus(session, orderId, constants.OrderSuccess)
	if err != nil {
		_ = session.Rollback()
		return nil, err
	}

	// Submit transaction.
	err = session.Commit()
	if err != nil {
		return nil, err
	}

	return &orderPb.SubmitOrderResponse{}, nil
}

func (s *Server) CancelOrder(ctx context.Context, in *orderPb.CancelOrderRequest) (*orderPb.CancelOrderResponse, error) {
	return nil, nil
}
