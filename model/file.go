package model

import (
	"github.com/TRON-US/soter-order-service/common/errorm"

	"github.com/go-xorm/xorm"
)

var (
	insertFileInfoSql = `INSERT INTO file (user_id, file_name, file_size, expire_time) VALUES (?, ?, ?, from_unixtime(?))`
	updateFileHashSql = `UPDATE file SET file_hash = ? WHERE id = ? AND file_hash IS NULL`
	deleteFileSql     = `UPDATE file SET deleted = 1 WHERE id = ? AND deleted = 0`
)

// Insert file info table.
func InsertFileInfo(session *xorm.Session, userId, fileSize int64, fileName string, expireTime int) (int64, error) {
	// Execute insert sql.
	r, err := session.Exec(insertFileInfoSql, userId, fileName, fileSize, expireTime)
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

// Update file hash by file id.
func UpdateFileHash(session *xorm.Session, id int64, fileHash string) error {
	// Execute update sql.
	r, err := session.Exec(updateFileHashSql, fileHash, id)
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

// Update file status to deleted.
func DeleteFile(session *xorm.Session, id int64) error {
	// Execute update sql.
	r, err := session.Exec(deleteFileSql, id)
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
