package model

import (
	"time"

	chaos "github.com/TRON-US/chaos/project/soter"
	"github.com/TRON-US/soter-order-service/common/errorm"

	"github.com/go-xorm/xorm"
)

var (
	queryLedgerInfoByAddressSql = `
			SELECT
				id,
				user_id,
				address,
				total_times,
				total_size,
				balance,
				freeze_balance,
				total_fee,
				unix_timestamp(update_time) as update_time,
				balance_check,
				version
			FROM
				ledger
			WHERE
				address = ?
		`
	insertLedgerInfoSql  = `INSERT INTO ledger (user_id, address, balance_check, update_time) VALUES (?, ?, ?, ?)`
	updateUserBalanceSql = `
			UPDATE
				ledger
			SET
				balance = ?,
				freeze_balance = ?,
				update_time = from_unixtime(?),
				balance_check = ?,
				version = version + 1
			WHERE
				id = ? AND version = ?
		`
	updateLedgerInfoSql = `
			UPDATE 
				ledger
			SET
				total_times = total_times + 1,
				total_size = ?,
				balance = ?,
				freeze_balance = ?,
				total_fee = ?,
				update_time = from_unixtime(?),
				balance_check = ?,
				version = version + 1
			WHERE
				id = ? AND version = ?
		`
)

type Ledger struct {
	Id            int64
	UserId        int64
	Address       string
	TotalTimes    int
	TotalSize     int64
	Balance       int64
	FreezeBalance int64
	TotalFee      int64
	UpdateTime    int
	BalanceCheck  string
	Version       int64
}

// Select ledger by user address.
func (db *Database) QueryLedgerInfoByAddress(address string) (*Ledger, error) {
	// Execute query sql.
	row := db.DB.DB().QueryRow(queryLedgerInfoByAddressSql, address)
	ledger := &Ledger{}
	err := row.Scan(&ledger.Id, &ledger.UserId, &ledger.Address, &ledger.TotalTimes, &ledger.TotalSize, &ledger.Balance, &ledger.FreezeBalance, &ledger.TotalFee, &ledger.UpdateTime, &ledger.BalanceCheck, &ledger.Version)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			// Open transaction.
			session := db.DB.NewSession()
			err := session.Begin()
			if err != nil {
				return nil, err
			}
			defer session.Close()

			// Init user info.
			userId, err := initUser(session, address)
			if err != nil {
				_ = session.Rollback()
				return nil, err
			}

			// Init ledger info
			ledger, err = initLedger(session, userId, address)
			if err != nil {
				_ = session.Rollback()
				return nil, err
			}

			// Commit transaction.
			err = session.Commit()
			if err != nil {
				return nil, err
			}

			return ledger, nil
		} else {
			return nil, err
		}
	}

	// Verify balance check.
	if !ledger.VerifyBalanceCheck() {
		return nil, errorm.AccountIllegal
	}
	return ledger, nil
}

// Init ledger information.
func initLedger(session *xorm.Session, userId int64, userAddress string) (*Ledger, error) {
	now := time.Now().Local()
	ledger := &Ledger{
		UserId:        userId,
		Address:       userAddress,
		Balance:       0,
		FreezeBalance: 0,
		UpdateTime:    int(now.Unix()),
	}

	// Get balance check.
	balanceCheck, err := ledger.GetBalanceCheck()
	if err != nil {
		return nil, err
	}

	// Execute SQL.
	r, err := session.Exec(insertLedgerInfoSql, userId, userAddress, balanceCheck, now.Format("2006-01-02 15:04:05"))
	if err != nil {
		return nil, err
	}

	// Get last insert id.
	id, err := r.LastInsertId()
	if err != nil {
		return nil, err
	}

	ledger.Id = id
	ledger.Version = 1

	return ledger, nil
}

// Update user balance and freeze balance by id and version.
func UpdateUserBalance(session *xorm.Session, balance, freezeBalance, version, id int64, address string, updateTime int) error {
	ledger := &Ledger{
		Address:       address,
		Balance:       balance,
		FreezeBalance: freezeBalance,
		UpdateTime:    updateTime,
	}

	// Get balance check
	balanceCheck, err := ledger.GetBalanceCheck()
	if err != nil {
		return err
	}

	// Execute update sql.
	r, err := session.Exec(updateUserBalanceSql, balance, freezeBalance, updateTime, balanceCheck, id, version)
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

// Update ledger information by id and version.
func UpdateLedgerInfo(session *xorm.Session, totalSize, balance, freezeBalance, totalFee, version, id int64, address string, updateTime int) error {
	ledger := &Ledger{
		Address:       address,
		Balance:       balance,
		FreezeBalance: freezeBalance,
		UpdateTime:    updateTime,
	}

	// Get balance check.
	balanceCheck, err := ledger.GetBalanceCheck()
	if err != nil {
		return err
	}

	// Execute update sql.
	r, err := session.Exec(updateLedgerInfoSql, totalSize, balance, freezeBalance, totalFee, updateTime, balanceCheck, id, version)
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

// Verify balance check.
func (ledger *Ledger) VerifyBalanceCheck() bool {
	balanceCheck := chaos.BalanceCheck{
		Address:       ledger.Address,
		Balance:       ledger.Balance,
		FreezeBalance: ledger.FreezeBalance,
		Timestamp:     ledger.UpdateTime,
	}

	return balanceCheck.VerifyBalanceCheck(ledger.BalanceCheck)
}

// Get balance check.
func (ledger *Ledger) GetBalanceCheck() (string, error) {
	balanceCheck := chaos.BalanceCheck{
		Address:       ledger.Address,
		Balance:       ledger.Balance,
		FreezeBalance: ledger.FreezeBalance,
		Timestamp:     ledger.UpdateTime,
	}

	return balanceCheck.GetBalanceCheck()
}
