package service

import (
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

func (s *Server) CreateOrder(ctx context.Context, in *orderPb.CreateOrderRequest) (*orderPb.CreateOrderResponse, error) {
	return nil, nil
}

func (s *Server) SubmitOrder(ctx context.Context, in *orderPb.SubmitOrderRequest) (*orderPb.SubmitOrderResponse, error) {
	return nil, nil
}

func (s *Server) CancelOrder(ctx context.Context, in *orderPb.CancelOrderRequest) (*orderPb.CancelOrderResponse, error) {
	return nil, nil
}
