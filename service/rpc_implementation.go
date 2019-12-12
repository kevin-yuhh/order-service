package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/TRON-US/soter-order-service/utils"

	"github.com/TRON-US/chaos/network/slack"
	"github.com/TRON-US/soter-order-service/charge"
	"github.com/TRON-US/soter-order-service/common/constants"
	"github.com/TRON-US/soter-order-service/common/errorm"
	"github.com/TRON-US/soter-order-service/config"
	"github.com/TRON-US/soter-order-service/logger"
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
		go func(address string, err error) {
			errMessage := fmt.Sprintf("Address: [%v] query ledger info error, reasons: [%v]", address, err)
			logger.Logger.Errorw(errMessage, "function", constants.QueryLedgerInfoModel)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.QueryLedgerInfoModel, constants.SlackNotifyLevel0),
				s.Config.Slack.SlackNotificationTimeout,
				slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
		}(address, err)
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
		go func(address, requestId string, err error) {
			errMessage := fmt.Sprintf("Address: [%v], requestId: [%v] query ledger info error, reasons: [%v]", address, requestId, err)
			logger.Logger.Errorw(errMessage, "function", constants.QueryLedgerInfoModel)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.QueryLedgerInfoModel, constants.SlackNotifyLevel0),
				s.Config.Slack.SlackNotificationTimeout,
				slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
		}(address, requestId, err)
		return nil, err
	}

	// Calculate fee of this order.
	fee := s.Fee.Fee(fileSize, ledger.TotalTimes, s.Time)

	// Get activity rate.
	rate, err := s.DbConn.QueryActivityByUserId(ledger.UserId)
	if err != nil {
		go func(address, requestId string, err error) {
			errMessage := fmt.Sprintf("Address: [%v], requestId: [%v] query activity info error, reasons: [%v]", address, requestId, err)
			logger.Logger.Errorw(errMessage, "function", constants.QueryActivityInfoModel)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.QueryActivityInfoModel, constants.SlackNotifyLevel0),
				s.Config.Slack.SlackNotificationTimeout,
				slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
		}(address, requestId, err)
		return nil, err
	}

	amount := int64(rate * float64(fee))

	// Check balance illegal.
	if ledger.Balance < amount {
		return nil, errorm.InsufficientBalance
	}

	// Open transaction.
	session := s.DbConn.DB.NewSession()
	err = session.Begin()
	if err != nil {
		go func(address, requestId string, err error) {
			errMessage := fmt.Sprintf("Address: [%v], requestId: [%v] open transaction error, reasons: [%v]", address, requestId, err)
			logger.Logger.Errorw(errMessage, "function", constants.SessionBegin)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.SessionBegin, constants.SlackNotifyLevel0),
				s.Config.Slack.SlackNotificationTimeout,
				slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
		}(address, requestId, err)
		return nil, err
	}
	defer session.Close()

	// Insert file information
	fileId, err := model.InsertFileInfo(session, ledger.UserId, fileSize, fileName, int(time.Now().Local().Unix())+s.Time*86400)
	if err != nil {
		_ = session.Rollback()
		go func(address, requestId string, err error) {
			errMessage := fmt.Sprintf("Address: [%v], requestId: [%v] insert file info error, reasons: [%v]", address, requestId, err)
			logger.Logger.Errorw(errMessage, "function", constants.InsertFileInfoModel)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.InsertFileInfoModel, constants.SlackNotifyLevel0),
				s.Config.Slack.SlackNotificationTimeout,
				slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
		}(address, requestId, err)
		return nil, err
	}

	// Insert order information.
	id, err := model.InsertOrderInfo(session, ledger.UserId, fileId, amount, s.Fee.StrategyId, requestId, s.Time)
	if err != nil {
		_ = session.Rollback()
		go func(address, requestId string, err error) {
			errMessage := fmt.Sprintf("Address: [%v], requestId: [%v] insert order info error, reasons: [%v]", address, requestId, err)
			logger.Logger.Errorw(errMessage, "function", constants.InsertOrderInfoModel)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.InsertOrderInfoModel, constants.SlackNotifyLevel0),
				s.Config.Slack.SlackNotificationTimeout,
				slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
		}(address, requestId, err)
		return nil, err
	}

	// Freeze user balance.
	err = model.UpdateUserBalance(session, ledger.Balance-amount, ledger.FreezeBalance+amount, ledger.Version, ledger.Id, ledger.Address, int(time.Now().Local().Unix()))
	if err != nil {
		_ = session.Rollback()
		go func(address, requestId string, err error) {
			errMessage := fmt.Sprintf("Address: [%v], requestId: [%v] update user balance error, reasons: [%v]", address, requestId, err)
			logger.Logger.Errorw(errMessage, "function", constants.UpdateUserBalanceModel)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.UpdateUserBalanceModel, constants.SlackNotifyLevel0),
				s.Config.Slack.SlackNotificationTimeout,
				slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
		}(address, requestId, err)
		return nil, err
	}

	// Submit transaction.
	err = session.Commit()
	if err != nil {
		go func(address, requestId string, err error) {
			errMessage := fmt.Sprintf("Address: [%v], requestId: [%v] commit transaction error, reasons: [%v]", address, requestId, err)
			logger.Logger.Errorw(errMessage, "function", constants.SessionCommit)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.SessionCommit, constants.SlackNotifyLevel0),
				s.Config.Slack.SlackNotificationTimeout,
				slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
		}(address, requestId, err)
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
		if err.Error() != errorm.QueryResultEmpty {
			go func(orderId int64, err error) {
				errMessage := fmt.Sprintf("orderId: [%v] query order info error, reasons: [%v]", orderId, err)
				logger.Logger.Errorw(errMessage, "function", constants.QueryOrderInfoModel)
				_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
					utils.ErrorRequestBody(errMessage, constants.QueryOrderInfoModel, constants.SlackNotifyLevel0),
					s.Config.Slack.SlackNotificationTimeout,
					slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
			}(orderId, err)
		}
		return nil, err
	}

	// Query ledger info by address.
	ledger, err := s.DbConn.QueryLedgerInfoByAddress(order.Address)
	if err != nil {
		go func(orderId int64, err error) {
			errMessage := fmt.Sprintf("orderId: [%v] query ledger info error, reasons: [%v]", orderId, err)
			logger.Logger.Errorw(errMessage, "function", constants.QueryLedgerInfoModel)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.QueryLedgerInfoModel, constants.SlackNotifyLevel0),
				s.Config.Slack.SlackNotificationTimeout,
				slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
		}(orderId, err)
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
		go func(orderId int64, err error) {
			errMessage := fmt.Sprintf("orderId: [%v] open transaction error, reasons: [%v]", orderId, err)
			logger.Logger.Errorw(errMessage, "function", constants.SessionBegin)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.SessionBegin, constants.SlackNotifyLevel0),
				s.Config.Slack.SlackNotificationTimeout,
				slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
		}(orderId, err)
		return nil, err
	}
	defer session.Close()

	// Update file hash by file id.
	err = model.UpdateFileHash(session, order.FileId, order.FileVersion, fileHash)
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062") {
			// Query file by user id and file hash.
			file, err := s.DbConn.QueryFileByUk(ledger.UserId, fileHash)
			if err != nil {
				_ = session.Rollback()
				go func(orderId int64, err error) {
					errMessage := fmt.Sprintf("orderId: [%v] query file by file hash error, reasons: [%v]", orderId, err)
					logger.Logger.Errorw(errMessage, "function", constants.QueryFileByUkModel)
					_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
						utils.ErrorRequestBody(errMessage, constants.QueryFileByUkModel, constants.SlackNotifyLevel0),
						s.Config.Slack.SlackNotificationTimeout,
						slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
				}(orderId, err)
				return nil, err
			}

			// Check new order expire time if greater than old order.
			if order.ExpireTime > file.ExpireTime {
				// Update old file expire time to new expire time.
				err = model.UpdateFileExpireTime(session, order.ExpireTime, file.Id, file.Version)
				if err != nil {
					_ = session.Rollback()
					go func(orderId int64, err error) {
						errMessage := fmt.Sprintf("orderId: [%v] update file expire time error, reasons: [%v]", orderId, err)
						logger.Logger.Errorw(errMessage, "function", constants.UpdateFileExpireTimeModel)
						_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
							utils.ErrorRequestBody(errMessage, constants.UpdateFileExpireTimeModel, constants.SlackNotifyLevel0),
							s.Config.Slack.SlackNotificationTimeout,
							slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
					}(orderId, err)
					return nil, err
				}
			}

			// Update order file id.
			err = model.UpdateOrderFileIdById(session, file.Id, orderId)
			if err != nil {
				_ = session.Rollback()
				go func(orderId int64, err error) {
					errMessage := fmt.Sprintf("orderId: [%v] update order info file id error, reasons: [%v]", orderId, err)
					logger.Logger.Errorw(errMessage, "function", constants.UpdateOrderFileIdModel)
					_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
						utils.ErrorRequestBody(errMessage, constants.UpdateOrderFileIdModel, constants.SlackNotifyLevel0),
						s.Config.Slack.SlackNotificationTimeout,
						slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
				}(orderId, err)
				return nil, err
			}

			// Delete file.
			err = model.DeleteFile(session, order.FileId, order.FileVersion)
			if err != nil {
				_ = session.Rollback()
				go func(orderId int64, err error) {
					errMessage := fmt.Sprintf("orderId: [%v] delete file error, reasons: [%v]", orderId, err)
					logger.Logger.Errorw(errMessage, "function", constants.DeleteFileModel)
					_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
						utils.ErrorRequestBody(errMessage, constants.DeleteFileModel, constants.SlackNotifyLevel0),
						s.Config.Slack.SlackNotificationTimeout,
						slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
				}(orderId, err)
				return nil, err
			}
		} else {
			_ = session.Rollback()
			go func(orderId int64, err error) {
				errMessage := fmt.Sprintf("orderId: [%v] update file hash error, reasons: [%v]", orderId, err)
				logger.Logger.Errorw(errMessage, "function", constants.UpdateFileHashModel)
				_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
					utils.ErrorRequestBody(errMessage, constants.UpdateFileHashModel, constants.SlackNotifyLevel0),
					s.Config.Slack.SlackNotificationTimeout,
					slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
			}(orderId, err)
			return nil, err
		}
	}

	// Update ledger information by ledger id.
	err = model.UpdateLedgerInfo(session, ledger.TotalSize+order.FileSize, ledger.Balance, ledger.FreezeBalance-order.Amount, ledger.TotalFee+order.Amount, ledger.Version, ledger.Id, order.Address, int(time.Now().Local().Unix()))
	if err != nil {
		_ = session.Rollback()
		go func(orderId int64, err error) {
			errMessage := fmt.Sprintf("orderId: [%v] update ledger info error, reasons: [%v]", orderId, err)
			logger.Logger.Errorw(errMessage, "function", constants.UpdateLedgerInfoModel)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.UpdateLedgerInfoModel, constants.SlackNotifyLevel0),
				s.Config.Slack.SlackNotificationTimeout,
				slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
		}(orderId, err)
		return nil, err
	}

	// Update order status by order id.
	err = model.UpdateOrderStatus(session, orderId, constants.OrderSuccess)
	if err != nil {
		_ = session.Rollback()
		go func(orderId int64, err error) {
			errMessage := fmt.Sprintf("orderId: [%v] update order status error, reasons: [%v]", orderId, err)
			logger.Logger.Errorw(errMessage, "function", constants.UpdateOrderStatusModel)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.UpdateOrderStatusModel, constants.SlackNotifyLevel0),
				s.Config.Slack.SlackNotificationTimeout,
				slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
		}(orderId, err)
		return nil, err
	}

	// Submit transaction.
	err = session.Commit()
	if err != nil {
		go func(orderId int64, err error) {
			errMessage := fmt.Sprintf("orderId: [%v] commit transaction error, reasons: [%v]", orderId, err)
			logger.Logger.Errorw(errMessage, "function", constants.SessionCommit)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.SessionCommit, constants.SlackNotifyLevel0),
				s.Config.Slack.SlackNotificationTimeout,
				slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
		}(orderId, err)
		return nil, err
	}

	return &orderPb.SubmitOrderResponse{}, nil
}

// Close order by order id.
func (s *Server) CloseOrder(ctx context.Context, in *orderPb.CloseOrderRequest) (*orderPb.CloseOrderResponse, error) {
	// Check input params.
	orderId := in.GetOrderId()
	if orderId <= 0 {
		return nil, errorm.RequestParamEmpty
	}

	// Get order info by order id.
	order, err := s.DbConn.QueryOrderInfoById(orderId)
	if err != nil {
		if err.Error() != errorm.QueryResultEmpty {
			go func(orderId int64, err error) {
				errMessage := fmt.Sprintf("orderId: [%v] query order info error, reasons: [%v]", orderId, err)
				logger.Logger.Errorw(errMessage, "function", constants.QueryOrderInfoModel)
				_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
					utils.ErrorRequestBody(errMessage, constants.QueryOrderInfoModel, constants.SlackNotifyLevel0),
					s.Config.Slack.SlackNotificationTimeout,
					slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
			}(orderId, err)
		}
		return nil, err
	}

	// Query ledger info by address.
	ledger, err := s.DbConn.QueryLedgerInfoByAddress(order.Address)
	if err != nil {
		go func(orderId int64, err error) {
			errMessage := fmt.Sprintf("orderId: [%v] query ledger info error, reasons: [%v]", orderId, err)
			logger.Logger.Errorw(errMessage, "function", constants.QueryLedgerInfoModel)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.QueryLedgerInfoModel, constants.SlackNotifyLevel0),
				s.Config.Slack.SlackNotificationTimeout,
				slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
		}(orderId, err)
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
		go func(orderId int64, err error) {
			errMessage := fmt.Sprintf("orderId: [%v] open transaction error, reasons: [%v]", orderId, err)
			logger.Logger.Errorw(errMessage, "function", constants.SessionBegin)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.SessionBegin, constants.SlackNotifyLevel0),
				s.Config.Slack.SlackNotificationTimeout,
				slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
		}(orderId, err)
		return nil, err
	}
	defer session.Close()

	// Delete file.
	err = model.DeleteFile(session, order.FileId, order.FileVersion)
	if err != nil {
		_ = session.Rollback()
		go func(orderId int64, err error) {
			errMessage := fmt.Sprintf("orderId: [%v] delete file error, reasons: [%v]", orderId, err)
			logger.Logger.Errorw(errMessage, "function", constants.DeleteFileModel)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.DeleteFileModel, constants.SlackNotifyLevel0),
				s.Config.Slack.SlackNotificationTimeout,
				slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
		}(orderId, err)
		return nil, err
	}

	// Update ledger information by ledger id.
	err = model.UpdateLedgerInfo(session, ledger.TotalSize, ledger.Balance+order.Amount, ledger.FreezeBalance-order.Amount, ledger.TotalFee, ledger.Version, ledger.Id, order.Address, int(time.Now().Local().Unix()))
	if err != nil {
		_ = session.Rollback()
		go func(orderId int64, err error) {
			errMessage := fmt.Sprintf("orderId: [%v] update ledger info error, reasons: [%v]", orderId, err)
			logger.Logger.Errorw(errMessage, "function", constants.UpdateLedgerInfoModel)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.UpdateLedgerInfoModel, constants.SlackNotifyLevel0),
				s.Config.Slack.SlackNotificationTimeout,
				slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
		}(orderId, err)
		return nil, err
	}

	// Update order status by order id.
	err = model.UpdateOrderStatus(session, orderId, constants.OrderFailed)
	if err != nil {
		_ = session.Rollback()
		go func(orderId int64, err error) {
			errMessage := fmt.Sprintf("orderId: [%v] update order status error, reasons: [%v]", orderId, err)
			logger.Logger.Errorw(errMessage, "function", constants.UpdateOrderStatusModel)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.UpdateOrderStatusModel, constants.SlackNotifyLevel0),
				s.Config.Slack.SlackNotificationTimeout,
				slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
		}(orderId, err)
		return nil, err
	}

	// Submit transaction.
	err = session.Commit()
	if err != nil {
		go func(orderId int64, err error) {
			errMessage := fmt.Sprintf("orderId: [%v] commit transaction error, reasons: [%v]", orderId, err)
			logger.Logger.Errorw(errMessage, "function", constants.SessionCommit)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.SessionCommit, constants.SlackNotifyLevel0),
				s.Config.Slack.SlackNotificationTimeout,
				slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
		}(orderId, err)
		return nil, err
	}

	return &orderPb.CloseOrderResponse{}, nil
}
