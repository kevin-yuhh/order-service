package service

import (
	"fmt"
	"time"

	"github.com/TRON-US/soter-order-service/common/constants"
	"github.com/TRON-US/soter-order-service/common/errorm"
	"github.com/TRON-US/soter-order-service/logger"
	"github.com/TRON-US/soter-order-service/model"
	"github.com/TRON-US/soter-order-service/utils"

	"github.com/TRON-US/chaos/network/slack"
)

// Create order controller.
func (s *Server) CreateOrderController(address, requestId, fileName string, fileSize int64) (*int64, error) {
	// Open transaction.
	session := s.DbConn.DB.NewSession()
	err := session.Begin()
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

	// Lock row by address.
	ledger, err := s.DbConn.LockLedgerInfoByAddress(session, address)
	if err != nil {
		_ = session.Rollback()
		go func(address, requestId string, err error) {
			errMessage := fmt.Sprintf("Address: [%v], requestId: [%v] get lock error, reasons: [%v]", address, requestId, err)
			logger.Logger.Errorw(errMessage, "function", constants.GetLedgerRowLockModel)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.GetLedgerRowLockModel, constants.SlackNotifyLevel0),
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
		_ = session.Rollback()
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
		_ = session.Rollback()
		return nil, errorm.InsufficientBalance
	}

	// Insert file information
	fileId, err := model.InsertFileInfo(session, ledger.UserId, fileSize, fileName, int(time.Now().Local().Unix())+s.Time*constants.DaySeconds)
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
	id, err := model.InsertOrderInfo(session, ledger.UserId, fileId, amount, s.Fee.StrategyId, requestId, constants.Charge, constants.OrderPending, s.Time)
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

	return &id, nil
}
