package model

var (
	queryConfigByEnvSql = `SELECT strategy_id, default_time FROM config WHERE env = ?`
)

// Select config by env.
func (db *Database) QueryConfig(env string) (int64, int, error) {
	var strategyId int64
	var time int

	// Execute query sql.
	row := db.DB.DB().QueryRow(queryConfigByEnvSql, env)
	err := row.Scan(&strategyId, &time)
	if err != nil {
		return 0, 0, err
	}
	return strategyId, time, nil
}
