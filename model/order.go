package model

import (
	"github.com/TRON-US/soter-order-service/common/errorm"

	"github.com/go-xorm/xorm"
)

var (
	insertOrderInfoSql = `INSERT INTO order_info (user_id, request_id, amount, strategy_id, time) VALUES (?,?,?,?,?)`
	updateOrderInfoSql = `UPDATE order_info SET file_id = ?, status = ? WHERE id = ? AND status = 'U'`
)

// Insert order info.
func InsertOrderInfo(session *xorm.Session, userId, amount, strategyId int64, requestId string, time int) (int64, error) {
	// Execute insert order info sql.
	r, err := session.Exec(insertOrderInfoSql, userId, requestId, amount, strategyId, time)
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

// Update order information by id and old status.
func UpdateOrderInfo(session *xorm.Session, fileId, id int64, status string) error {
	// Execute update order info sql.
	r, err := session.Exec(updateOrderInfoSql, fileId, status, id)
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
