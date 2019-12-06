package model

import (
	"github.com/go-xorm/xorm"
)

var (
	initUserSql = `INSERT INTO user (address) VALUES (?)`
)

// Init user information,
func initUser(session *xorm.Session, userAddress string) (int64, error) {
	// Execute SQL.
	r, err := session.Exec(initUserSql, userAddress)
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
