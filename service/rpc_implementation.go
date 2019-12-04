package service

import (
	"fmt"
	"time"

	"github.com/TRON-US/soter-order-service/charge"
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

	// Call QueryLedgerInfoByAddress api.
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

	// Call QueryLedgerInfoByAddress api.
	ledger, err := s.DbConn.QueryLedgerInfoByAddress(address)
	if err != nil {
		return nil, err
	}

	fmt.Println(fileSize, ledger.TotalTimes, s.Time)

	// Calculate fee of this order.
	amount := s.Fee.Fee(fileSize, ledger.TotalTimes, s.Time)

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

func (s *Server) SubmitOrder(ctx context.Context, in *orderPb.SubmitOrderRequest) (*orderPb.SubmitOrderResponse, error) {
	return nil, nil
}

func (s *Server) CancelOrder(ctx context.Context, in *orderPb.CancelOrderRequest) (*orderPb.CancelOrderResponse, error) {
	return nil, nil
}
