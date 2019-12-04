package model

import (
	"github.com/TRON-US/soter-order-service/common/errorm"

	"github.com/go-xorm/xorm"
)

var (
	queryFileByUkSql        = `SELECT id, file_name, file_size, unix_timestamp(expire_time), deleted, version FROM file WHERE user_id = ? AND file_hash = ?`
	insertFileInfoSql       = `INSERT INTO file (user_id, file_name, file_size, expire_time) VALUES (?, ?, ?, from_unixtime(?))`
	updateFileHashSql       = `UPDATE file SET file_hash = ?, version = version + 1 WHERE id = ? AND version = ? AND file_hash IS NULL`
	updateFileExpireTimeSql = `UPDATE file SET expire_time = from_unixtime(?), version = version + 1 WHERE id = ? AND version = ?`
	deleteFileSql           = `UPDATE file SET deleted = 1, version = version + 1 WHERE id = ? AND deleted = 0 AND version = ?`
)

type File struct {
	Id         int64
	FileName   string
	FileSize   int64
	ExpireTime int64
	Deleted    int8
	Version    int64
}

// Select file information by UK.
func (db *Database) QueryFileByUk(userId int64, fileHash string) (*File, error) {
	// Execute query sql.
	row := db.DB.DB().QueryRow(queryFileByUkSql, userId, fileHash)
	file := &File{}
	err := row.Scan(&file.Id, &file.FileName, &file.FileSize, &file.ExpireTime, &file.Deleted, &file.Version)
	if err != nil {
		return nil, err
	}

	return file, nil
}

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
func UpdateFileHash(session *xorm.Session, id, version int64, fileHash string) error {
	// Execute update sql.
	r, err := session.Exec(updateFileHashSql, fileHash, id, version)
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

// Update file expire time by file Id.
func UpdateFileExpireTime(session *xorm.Session, expireTime, id, version int64) error {
	// Execute update sql.
	r, err := session.Exec(updateFileExpireTimeSql, expireTime, id, version)
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
func DeleteFile(session *xorm.Session, id, version int64) error {
	// Execute update sql.
	r, err := session.Exec(deleteFileSql, id, version)
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
