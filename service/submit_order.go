package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/TRON-US/soter-order-service/common/constants"
	"github.com/TRON-US/soter-order-service/common/errorm"
	"github.com/TRON-US/soter-order-service/logger"
	"github.com/TRON-US/soter-order-service/model"
	"github.com/TRON-US/soter-order-service/utils"

	"github.com/Shopify/sarama"
	"github.com/TRON-US/chaos/network/slack"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/go-xorm/xorm"
	"github.com/prometheus/client_golang/prometheus"
)

type OrderNotify struct {
	Result    string `json:"result"`
	OrderId   int64  `json:"order_id"`
	SessionId string `json:"session_id"`
	FileHash  string `json:"file_hash"`
}

// Process success or error order info.
func (s *Server) ClusterConsumer() {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

	// TODO Read processing offset from redis.
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// Init kafka consumer.
	consumer, err := cluster.NewConsumer(s.Config.Kafka.Servers, s.Config.Kafka.GroupId, s.Config.Kafka.Topic, config)
	if err != nil {
		logger.Logger.Fatalf("GroupId: [%v] new consumer error, reasons: [%v]", s.Config.Kafka.GroupId, err)
	}
	defer consumer.Close()

	// Consume errors.
	go func() {
		for err := range consumer.Errors() {
			logger.Logger.Errorf("GroupId: [%v] consume error, reasons: [%v]", s.Config.Kafka.GroupId, err)
		}
	}()

	// Consume notifications.
	go func() {
		for ntf := range consumer.Notifications() {
			logger.Logger.Infof("%s:Rebalanced: %+v", s.Config.Kafka.GroupId, ntf)
		}
	}()

	// Consume messages.
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				t := time.Now()
				rpcRequestCount.With(prometheus.Labels{"method": "SubmitOrderTotal"}).Inc()

				order := OrderNotify{}
				err := json.Unmarshal(msg.Value, &order)
				if err != nil {
					rpcRequestCount.With(prometheus.Labels{"method": "SubmitOrderFailed"}).Inc()
					consumer.MarkOffset(msg, constants.KafkaResultError)
					continue
				}

				// Submit order by file hash, result and order id.
				err = s.SubmitOrderController(order.FileHash, order.Result, order.OrderId)
				if err != nil {
					rpcRequestCount.With(prometheus.Labels{"method": "SubmitOrderError"}).Inc()
					rpcRequestDuration.With(prometheus.Labels{"method": "SubmitOrder"}).Observe(float64(time.Since(t).Microseconds()) / 1000)
					consumer.MarkOffset(msg, constants.KafkaResultFailed)
					continue
				}

				rpcRequestCount.With(prometheus.Labels{"method": "SubmitOrderSuccess"}).Inc()
				rpcRequestDuration.With(prometheus.Labels{"method": "SubmitOrder"}).Observe(float64(time.Since(t).Microseconds()) / 1000)
				consumer.MarkOffset(msg, constants.KafkaResultSuccess)
			}
		}
	}
}

// Submit order controller.
func (s *Server) SubmitOrderController(fileHash, result string, orderId int64) error {
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

	// Check order type.
	if order.OrderType != constants.Charge {
		return errorm.OrderTypeIllegal
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

	if result == constants.BtfsNodeAgentError {
		// Query if exists file on same file hash.
		btfsFile, err := s.DbConn.QueryBtfsFileByHash(fileHash)
		if err != nil {
			if err.Error() != errorm.QueryResultEmpty {
				_ = session.Rollback()
				go func(orderId int64, err error) {
					errMessage := fmt.Sprintf("orderId: [%v] query btfs file by hash error, reasons: [%v]", orderId, err)
					logger.Logger.Errorw(errMessage, "function", constants.QueryBtfsFileByHashModel)
					_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
						utils.ErrorRequestBody(errMessage, constants.QueryBtfsFileByHashModel, constants.SlackNotifyLevel0),
						s.Config.Slack.SlackNotificationTimeout,
						slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
				}(orderId, err)
				return err
			}
		}

		// Check btfs_file if expired.
		if btfsFile == nil || time.Now().Local().Unix() > btfsFile.ExpireTime {
			// Refund.
			err = s.refund(session, *order, *ledger)
			if err != nil {
				_ = session.Rollback()
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

		// Update btfs file id.
		err = model.UpdateBtfsFileId(session, btfsFile.Id, order.FileId, order.FileVersion)
		if err != nil {
			if strings.Contains(err.Error(), errorm.DuplicateKey) {
				// Query file info by UK.
				file, err := s.DbConn.QueryFileByUk(order.UserId, btfsFile.Id)
				if err != nil {
					_ = session.Rollback()
					go func(orderId int64, err error) {
						errMessage := fmt.Sprintf("orderId: [%v] query file info by uk error, reasons: [%v]", orderId, err)
						logger.Logger.Errorw(errMessage, "function", constants.QueryFileByUkModel)
						_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
							utils.ErrorRequestBody(errMessage, constants.QueryFileByUkModel, constants.SlackNotifyLevel0),
							s.Config.Slack.SlackNotificationTimeout,
							slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
					}(orderId, err)
					return err
				}

				if time.Now().Local().Unix() <= file.ExpireTime {
					// Refund.
					err = s.refund(session, *order, *ledger)
					if err != nil {
						_ = session.Rollback()
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
				} else {
					// Reopen file by file id.
					err = model.ReopenFile(session, file.Id, file.Version, order.ExpireTime)
					if err != nil {
						_ = session.Rollback()
						go func(orderId int64, err error) {
							errMessage := fmt.Sprintf("orderId: [%v] reopen file error, reasons: [%v]", orderId, err)
							logger.Logger.Errorw(errMessage, "function", constants.ReopenFileModel)
							_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
								utils.ErrorRequestBody(errMessage, constants.ReopenFileModel, constants.SlackNotifyLevel0),
								s.Config.Slack.SlackNotificationTimeout,
								slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
						}(orderId, err)
						return err
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
						return err
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
						return err
					}
				}
			} else {
				_ = session.Rollback()
				go func(orderId int64, err error) {
					errMessage := fmt.Sprintf("orderId: [%v] update btfs file id error, reasons: [%v]", orderId, err)
					logger.Logger.Errorw(errMessage, "function", constants.UpdateBtfsFileIdModel)
					_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
						utils.ErrorRequestBody(errMessage, constants.UpdateBtfsFileIdModel, constants.SlackNotifyLevel0),
						s.Config.Slack.SlackNotificationTimeout,
						slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
				}(orderId, err)
				return err
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
			return err
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
	} else if result == constants.BtfsNodeAgentComplete {
		// Query if exists file on same file hash.
		btfsFile, err := s.DbConn.QueryBtfsFileByHash(fileHash)
		if err != nil {
			if err.Error() == errorm.QueryResultEmpty {
				// Insert into btfs file info
				btfsFileId, err := model.InsertBtfsFileInfo(session, fileHash, order.ExpireTime)
				if err != nil {
					_ = session.Rollback()
					go func(orderId int64, err error) {
						errMessage := fmt.Sprintf("orderId: [%v] insert btfs file info error, reasons: [%v]", orderId, err)
						logger.Logger.Errorw(errMessage, "function", constants.InsertBtfsFileInfoModel)
						_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
							utils.ErrorRequestBody(errMessage, constants.InsertBtfsFileInfoModel, constants.SlackNotifyLevel0),
							s.Config.Slack.SlackNotificationTimeout,
							slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
					}(orderId, err)
					return err
				}

				btfsFile = &model.BtfsFile{
					Id:         btfsFileId,
					FileHash:   fileHash,
					ExpireTime: order.ExpireTime,
					Version:    1,
				}
			} else {
				_ = session.Rollback()
				go func(orderId int64, err error) {
					errMessage := fmt.Sprintf("orderId: [%v] query btfs file by hash error, reasons: [%v]", orderId, err)
					logger.Logger.Errorw(errMessage, "function", constants.QueryBtfsFileByHashModel)
					_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
						utils.ErrorRequestBody(errMessage, constants.QueryBtfsFileByHashModel, constants.SlackNotifyLevel0),
						s.Config.Slack.SlackNotificationTimeout,
						slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
				}(orderId, err)
				return err
			}
		}

		if order.ExpireTime > btfsFile.ExpireTime {
			// Update btfs file expire time.
			err = model.UpdateBtfsFileExpireTime(session, btfsFile.Id, btfsFile.Version, order.ExpireTime)
			if err != nil {
				_ = session.Rollback()
				go func(orderId int64, err error) {
					errMessage := fmt.Sprintf("orderId: [%v] update btfs file expire time error, reasons: [%v]", orderId, err)
					logger.Logger.Errorw(errMessage, "function", constants.UpdateBtfsFileExpireTimeModel)
					_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
						utils.ErrorRequestBody(errMessage, constants.UpdateBtfsFileExpireTimeModel, constants.SlackNotifyLevel0),
						s.Config.Slack.SlackNotificationTimeout,
						slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
				}(orderId, err)
				return err
			}
		}

		// Update btfs file id.
		err = model.UpdateBtfsFileId(session, btfsFile.Id, order.FileId, order.FileVersion)
		if err != nil {
			if strings.Contains(err.Error(), errorm.DuplicateKey) {
				// Query file by uk.
				file, err := s.DbConn.QueryFileByUk(order.UserId, btfsFile.Id)
				if err != nil {
					_ = session.Rollback()
					go func(orderId int64, err error) {
						errMessage := fmt.Sprintf("orderId: [%v] query file info by uk error, reasons: [%v]", orderId, err)
						logger.Logger.Errorw(errMessage, "function", constants.QueryFileByUkModel)
						_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
							utils.ErrorRequestBody(errMessage, constants.QueryFileByUkModel, constants.SlackNotifyLevel0),
							s.Config.Slack.SlackNotificationTimeout,
							slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
					}(orderId, err)
					return err
				}

				// Reopen file by file id.
				err = model.ReopenFile(session, file.Id, file.Version, order.ExpireTime)
				if err != nil {
					_ = session.Rollback()
					go func(orderId int64, err error) {
						errMessage := fmt.Sprintf("orderId: [%v] reopen file error, reasons: [%v]", orderId, err)
						logger.Logger.Errorw(errMessage, "function", constants.ReopenFileModel)
						_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
							utils.ErrorRequestBody(errMessage, constants.ReopenFileModel, constants.SlackNotifyLevel0),
							s.Config.Slack.SlackNotificationTimeout,
							slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
					}(orderId, err)
					return err
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
					return err
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
					return err
				}
			} else {
				_ = session.Rollback()
				go func(orderId int64, err error) {
					errMessage := fmt.Sprintf("orderId: [%v] update btfs file id error, reasons: [%v]", orderId, err)
					logger.Logger.Errorw(errMessage, "function", constants.UpdateBtfsFileIdModel)
					_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
						utils.ErrorRequestBody(errMessage, constants.UpdateBtfsFileIdModel, constants.SlackNotifyLevel0),
						s.Config.Slack.SlackNotificationTimeout,
						slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
				}(orderId, err)
				return err
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
			return err
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
	} else {
		return errorm.BtfsStatusNotExists
	}
}

// Refund function.
func (s *Server) refund(session *xorm.Session, order model.Order, ledger model.Ledger) error {
	// Unfreeze user balance.
	err := model.UpdateUserBalance(session, ledger.Balance+order.Amount, ledger.FreezeBalance-order.Amount, ledger.Version, ledger.Id, order.Address, int(time.Now().Local().Unix()))
	if err != nil {
		go func(orderId int64, err error) {
			errMessage := fmt.Sprintf("orderId: [%v] update user balance error, reasons: [%v]", orderId, err)
			logger.Logger.Errorw(errMessage, "function", constants.UpdateUserBalanceModel)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.UpdateUserBalanceModel, constants.SlackNotifyLevel0),
				s.Config.Slack.SlackNotificationTimeout,
				slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
		}(order.Id, err)
		return err
	}

	// Update order status by order id.
	err = model.UpdateOrderStatus(session, order.Id, constants.OrderFailed)
	if err != nil {
		go func(orderId int64, err error) {
			errMessage := fmt.Sprintf("orderId: [%v] update order status error, reasons: [%v]", orderId, err)
			logger.Logger.Errorw(errMessage, "function", constants.UpdateOrderStatusModel)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.UpdateOrderStatusModel, constants.SlackNotifyLevel0),
				s.Config.Slack.SlackNotificationTimeout,
				slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
		}(order.Id, err)
		return err
	}

	// Delete file.
	err = model.DeleteFile(session, order.FileId, order.FileVersion)
	if err != nil {
		go func(orderId int64, err error) {
			errMessage := fmt.Sprintf("orderId: [%v] delete file error, reasons: [%v]", orderId, err)
			logger.Logger.Errorw(errMessage, "function", constants.DeleteFileModel)
			_ = slack.SendSlackNotification(s.Config.Slack.SlackWebhookUrl,
				utils.ErrorRequestBody(errMessage, constants.DeleteFileModel, constants.SlackNotifyLevel0),
				s.Config.Slack.SlackNotificationTimeout,
				slack.Priority0, slack.Priority(s.Config.Slack.SlackPriorityThreshold))
		}(order.Id, err)
		return err
	}

	return nil
}
