package model

import "github.com/go-xorm/xorm"

var (
	insertFileInfoSql = `INSERT INTO file (user_id, file_hash, file_name, file_size, expire_time) VALUES (?, ?, ?, ?, from_unixtime(?))`
)

// Insert file info table.
func InsertFileInfo(session *xorm.Session, userId, fileSize int64, fileHash, fileName string, expireTime int) (int64, error) {
	// Execute insert sql.
	r, err := session.Exec(insertFileInfoSql, userId, fileHash, fileName, fileSize, expireTime)
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
