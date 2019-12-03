package service

import (
	"time"

	"github.com/TRON-US/soter-order-service/common/errorm"
	"github.com/TRON-US/soter-order-service/model"
	orderPb "github.com/TRON-US/soter-order-service/proto"

	"golang.org/x/net/context"
)

type Server struct {
	DbConn *model.Database
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
	amount := in.GetAmount()
	if address == "" || requestId == "" || amount <= 0 {
		return nil, errorm.RequestParamEmpty
	}

	// Call QueryLedgerInfoByAddress api.
	ledger, err := s.DbConn.QueryLedgerInfoByAddress(address)
	if err != nil {
		return nil, err
	}

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

	// Insert order information.
	id, err := model.InsertOrderInfo(session, ledger.UserId, amount, s.DbConn.Config.ScriptId, requestId, s.DbConn.Config.DefaultTime)
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
