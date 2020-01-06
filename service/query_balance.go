package service

import (
	"fmt"

	"github.com/TRON-US/soter-order-service/common/constants"
	"github.com/TRON-US/soter-order-service/logger"
	"github.com/TRON-US/soter-order-service/model"
	"github.com/TRON-US/soter-order-service/utils"

	"github.com/TRON-US/chaos/network/slack"
)

// Query balance controller.
func (s *Server) QueryBalanceController(address string) (*model.Ledger, error) {
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
	return ledger, nil
}
