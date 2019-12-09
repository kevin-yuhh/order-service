package model

var (
	queryActivityByUserIdSql = `
		SELECT 
			b.strategy 
		FROM 
			activity_user a 
		LEFT JOIN 
			activity b 
		ON 
			a.activity_id = b.id 
		WHERE 
			a.user_id = ? 
		AND 
			b.begin_time <= now() 
		AND 
			b.end_time >= now() 
		AND 
			status = 0
	`
)

// Select activity strategy from activity.
func (db *Database) QueryActivityByUserId(userId int64) (float64, error) {
	var strategy float64

	// Execute query sql.
	row := db.DB.DB().QueryRow(queryActivityByUserIdSql, userId)
	err := row.Scan(&strategy)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return 1, nil
		} else {
			return 0, err
		}
	}

	return strategy, nil
}
