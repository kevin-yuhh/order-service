package model

import (
	"github.com/TRON-US/soter-order-service/common/errorm"

	"github.com/go-xorm/xorm"
)

var (
	queryBtfsFileInfoByHash     = `SELECT id, file_hash, unix_timestamp(expire_time), version FROM btfs_file WHERE file_hash = ?`
	insertBtfsFileSql           = `INSERT INTO btfs_file (file_hash, expire_time) VALUES (?, from_unixtime(?))`
	updateBtfsFileExpireTimeSql = `UPDATE btfs_file SET expire_time = from_unixtime(?), version = version + 1 WHERE id = ? AND version = ?`
)

type BtfsFile struct {
	Id         int64
	FileHash   string
	ExpireTime int64
	Version    int64
}

// Insert btfs file info table.
func InsertBtfsFileInfo(session *xorm.Session, fileHash string, expireTime int64) (int64, error) {
	// Execute insert sql.
	r, err := session.Exec(insertBtfsFileSql, fileHash, expireTime)
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

// Select btfs file information by file hash.
func (db *Database) QueryBtfsFileByHash(fileHash string) (*BtfsFile, error) {
	// Execute query sql.
	row := db.DB.DB().QueryRow(queryBtfsFileInfoByHash, fileHash)
	file := &BtfsFile{}
	err := row.Scan(&file.Id, &file.FileHash, &file.ExpireTime, &file.Version)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// Update btfs file expire time by file id.
func UpdateBtfsFileExpireTime(session *xorm.Session, id, version, expireTime int64) error {
	// Execute update sql.
	r, err := session.Exec(updateBtfsFileExpireTimeSql, expireTime, id, version)
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
