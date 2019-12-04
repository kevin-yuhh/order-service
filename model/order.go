package model

import (
	"github.com/TRON-US/soter-order-service/common/errorm"

	"github.com/go-xorm/xorm"
)

var (
	insertOrderInfoSql = `INSERT INTO order_info (user_id, file_id, request_id, amount, strategy_id, time, status) VALUES (?,?,?,?,?,?,'P')`
	updateOrderInfoSql = `UPDATE order_info SET status = ? WHERE id = ? AND status = 'P'`
)

// Insert order info.
func InsertOrderInfo(session *xorm.Session, userId, fileId, amount, strategyId int64, requestId string, time int) (int64, error) {
	// Execute insert order info sql.
	r, err := session.Exec(insertOrderInfoSql, userId, fileId, requestId, amount, strategyId, time)
	if err != nil {
		return 0, err
	}

	// Get last insert id.
	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update order status by id and old status.
func UpdateOrderStatus(session *xorm.Session, id int64, status string) error {
	// Execute update order info sql.
	r, err := session.Exec(updateOrderInfoSql, status, id)
	if err != nil {
		return err
	}

	// Get affected number.
	affected, err := r.RowsAffected()
	if err != nil {
		return err
	}

	// Row has not changed.
	if affected != 1 {
		return errorm.RowNotChanged
	}
	return nil
}
