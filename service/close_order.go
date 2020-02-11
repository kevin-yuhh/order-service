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

// Close order controller.
func (s *Server) CloseOrderController(orderId int64) error {
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
		return err
	}

	// Check order status.
	if order.Status != constants.OrderPending {
		return errorm.OrderStatusIllegal
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
		return err
	}
	defer session.Close()

	// Lock row by address.
	ledger, err := s.DbConn.LockLedgerInfoByAddress(session, order.Address)
	if err != nil {
		_ = session.Rollback()
		go func(orderId int64, err error) {
			errMessage := fmt.Sprintf("orderId: [%v] get lock error, reasons: [%v]", orderId, err)
			logger.Logger.Errorw(errMessage, "function", constants.GetLedgerRowLockModel)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.GetLedgerRowLockModel, constants.SlackNotifyLevel0),
				s.Config.Slack.SlackNotificationTimeout,
				slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
		}(orderId, err)
		return err
	}

	// Check freeze balance illegal.
	if ledger.FreezeBalance < order.Amount {
		_ = session.Rollback()
		return errorm.InsufficientBalance
	}

	if order.OrderType == constants.Charge {
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
			return err
		}
	}

	// Unfreeze user balance.
	err = model.UpdateUserBalance(session, ledger.Balance+order.Amount, ledger.FreezeBalance-order.Amount, ledger.Version, ledger.Id, order.Address, int(time.Now().Local().Unix()))
	if err != nil {
		_ = session.Rollback()
		go func(orderId int64, err error) {
			errMessage := fmt.Sprintf("orderId: [%v] update user balance error, reasons: [%v]", orderId, err)
			logger.Logger.Errorw(errMessage, "function", constants.UpdateUserBalanceModel)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.UpdateUserBalanceModel, constants.SlackNotifyLevel0),
				s.Config.Slack.SlackNotificationTimeout,
				slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
		}(orderId, err)
		return err
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
		return err
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
		return err
	}

	return nil
}
