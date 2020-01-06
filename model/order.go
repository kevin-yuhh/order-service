package model

import (
	"github.com/TRON-US/soter-order-service/common/errorm"

	"github.com/go-xorm/xorm"
)

var (
	insertOrderInfoSql        = `INSERT INTO order_info (user_id, file_id, type, request_id, amount, strategy_id, time, status) VALUES (?,?,?,?,?,?,?,?)`
	updateOrderInfoSql        = `UPDATE order_info SET status = ? WHERE id = ? AND status = 'P'`
	updateOrderFileIdByIdSql  = `UPDATE order_info SET file_id = ? WHERE id = ? AND status = 'P'`
	updateOrderSessionByIdSql = `UPDATE order_info SET session_id = ?, btfs_ip = ? WHERE id = ? AND status = 'P'`
	queryOrderInfoById        = `
		SELECT
			o.id,
			o.user_id, 
			o.file_id, 
			o.type, 
			o.request_id, 
			o.time,
			o.amount, 
			o.status, 
			u.address, 
			f.file_size,
			f.file_name,
			IFNULL(f.file_hash,''),
			unix_timestamp(f.expire_time),
			f.version
		FROM 
			order_info o 
		LEFT JOIN 
			user u 
		ON 
			o.user_id = u.id
		LEFT JOIN 
			(
				SELECT 
					a.id,
					a.file_size, 
					a.file_name, 
					b.file_hash, 
					a.expire_time,
					a.version 
				FROM 
					file a
				LEFT JOIN
					btfs_file b
				ON
					a.btfs_file_id = b.id
			) f
		ON
			o.file_id = f.id 
		WHERE 
			o.id = ?
		`
	queryOrderInfoRequestId = `
		SELECT
			o.id,
			o.user_id, 
			o.file_id, 
			o.type, 
			o.request_id, 
			o.time,
			o.amount, 
			o.status, 
			u.address, 
			f.file_size,
			f.file_name,
			IFNULL(f.file_hash,''),
			unix_timestamp(f.expire_time),
			f.version
		FROM 
			order_info o 
		LEFT JOIN 
			user u 
		ON 
			o.user_id = u.id
		LEFT JOIN 
			(
				SELECT 
					a.id,
					a.file_size, 
					a.file_name, 
					b.file_hash, 
					a.expire_time,
					a.version 
				FROM 
					file a
				LEFT JOIN
					btfs_file b
				ON
					a.btfs_file_id = b.id
			) f
		ON
			o.file_id = f.id 
		WHERE 
			o.request_id = ?
		AND
			u.address = ?
		`
)

type Order struct {
	Id          int64
	UserId      int64
	FileId      int64
	OrderType   string
	RequestId   string
	Times       int
	Amount      int64
	Status      string
	Address     string
	FileSize    int64
	FileName    string
	FileHash    string
	ExpireTime  int64
	FileVersion int64
}

// Insert order info.
func InsertOrderInfo(session *xorm.Session, userId, fileId, amount, strategyId int64, requestId, orderType, status string, time int) (int64, error) {
	// Execute insert order info sql.
	r, err := session.Exec(insertOrderInfoSql, userId, fileId, orderType, requestId, amount, strategyId, time, status)
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

// Update order file id by order id.
func UpdateOrderFileIdById(session *xorm.Session, fileId, id int64) error {
	// Execute update order info sql.
	r, err := session.Exec(updateOrderFileIdByIdSql, fileId, id)
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

// Update order session id and btfs ip address by id.
func UpdateOrderSessionById(session *xorm.Session, sessionId, nodeIp string, id int64) error {
	// Execute update order info sql.
	_, err := session.Exec(updateOrderSessionByIdSql, sessionId, nodeIp, id)
	if err != nil {
		return err
	}
	return nil
}

// Query order info by order id.
func (db *Database) QueryOrderInfoById(id int64) (*Order, error) {
	// Execute query sql.
	row := db.DB.DB().QueryRow(queryOrderInfoById, id)
	order := &Order{}
	err := row.Scan(&order.Id, &order.UserId, &order.FileId, &order.OrderType, &order.RequestId, &order.Times, &order.Amount, &order.Status, &order.Address, &order.FileSize, &order.FileName, &order.FileHash, &order.ExpireTime, &order.FileVersion)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// Query order info by request id and address.
func (db *Database) QueryOrderInfoByRequestId(requestId, address string) (*Order, error) {
	// Execute query sql.
	row := db.DB.DB().QueryRow(queryOrderInfoRequestId, requestId, address)
	order := &Order{}
	err := row.Scan(&order.Id, &order.UserId, &order.FileId, &order.OrderType, &order.RequestId, &order.Times, &order.Amount, &order.Status, &order.Address, &order.FileSize, &order.FileName, &order.FileHash, &order.ExpireTime, &order.FileVersion)
	if err != nil {
		if err.Error() == errorm.QueryResultEmpty {
			return nil, nil
		}
		return nil, err
	}

	return order, nil
}
