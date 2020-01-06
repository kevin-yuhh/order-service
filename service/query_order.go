package service

import (
	"fmt"

	"github.com/TRON-US/soter-order-service/common/constants"
	"github.com/TRON-US/soter-order-service/common/errorm"
	"github.com/TRON-US/soter-order-service/logger"
	"github.com/TRON-US/soter-order-service/model"
	"github.com/TRON-US/soter-order-service/utils"

	"github.com/TRON-US/chaos/network/slack"
)

// Query order controller.
func (s *Server) QueryOrderController(requestId, address string) (*model.Order, error) {
	// Query order info by request id and address.
	order, err := s.DbConn.QueryOrderInfoByRequestId(requestId, address)
	if err != nil {
		go func(address, requestId string, err error) {
			errMessage := fmt.Sprintf("Address: [%v], requestId: [%v] query order info error, reasons: [%v]", address, requestId, err)
			logger.Logger.Errorw(errMessage, "function", constants.QueryOrderInfo1Model)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.QueryOrderInfo1Model, constants.SlackNotifyLevel0),
				s.Config.Slack.SlackNotificationTimeout,
				slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
		}(address, requestId, err)
		return nil, err
	}

	// Check order info is nil.
	if order == nil {
		return nil, errorm.OrderNotExists
	}

	return order, nil
}
