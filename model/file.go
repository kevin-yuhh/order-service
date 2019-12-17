package model

import (
	"github.com/TRON-US/soter-order-service/common/errorm"

	"github.com/go-xorm/xorm"
)

var (
	queryFileByIdSql = `
		SELECT
			b.address, 
			a.file_name, 
			a.file_size, 
			unix_timestamp(a.expire_time), 
			a.file_hash,
			a.deleted, 
			a.version 
		FROM 
			file a 
		LEFT JOIN
			user b
		ON
			a.user_id = b.id
		WHERE
			a.id = ?
		`
	queryFileByUkSql        = `SELECT id, file_name, file_size, unix_timestamp(expire_time), deleted, version FROM file WHERE user_id = ? AND file_hash = ?`
	queryMaxExpireByHashSql = `select IFNULL(unix_timestamp(max(expire_time)),0) as expire_time from file where file_hash = ?`
	insertFileInfoSql       = `INSERT INTO file (user_id, file_name, file_size, expire_time) VALUES (?, ?, ?, from_unixtime(?))`
	updateFileHashSql       = `UPDATE file SET file_hash = ?, version = version + 1 WHERE id = ? AND version = ? AND file_hash IS NULL`
	updateFileExpireTimeSql = `UPDATE file SET expire_time = from_unixtime(?), version = version + 1 WHERE id = ? AND version = ?`
	deleteFileSql           = `UPDATE file SET deleted = 1, version = version + 1 WHERE id = ? AND deleted = 0 AND version = ?`
)

type File struct {
	Address    string
	Id         int64
	FileName   string
	FileSize   int64
	FileHash   string
	ExpireTime int64
	Deleted    int8
	Version    int64
}

// Select file information by id.
func (db *Database) QueryFileById(id int64) (*File, error) {
	// Execute query sql.
	row := db.DB.DB().QueryRow(queryFileByIdSql, id)
	file := &File{}
	err := row.Scan(&file.Address, &file.FileName, &file.FileSize, &file.ExpireTime, &file.FileHash, &file.Deleted, &file.Version)
	if err != nil {
		return nil, err
	}

	file.Id = id

	return file, nil
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

// Select max expire time by file hash.
func (db *Database) QueryMaxExpireByHash(fileHash string) (int64, error) {
	// Execute query sql.
	row := db.DB.DB().QueryRow(queryMaxExpireByHashSql, fileHash)

	var maxExpireTime int64

	err := row.Scan(&maxExpireTime)
	if err != nil {
		return 0, err
	}

	return maxExpireTime, nil
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
