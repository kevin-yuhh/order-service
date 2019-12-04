package model

var (
	queryStrategyByIdSql = `SELECT type, lua_script FROM strategy WHERE id = ?`
)

type Strategy struct {
	StrategyType int8
	Script       string
}

// Query fee strategy by id.
func (db *Database) QueryStrategyById(id int64) (*Strategy, error) {
	strategy := &Strategy{}

	// Execute query sql.
	row := db.DB.DB().QueryRow(queryStrategyByIdSql, id)
	err := row.Scan(&strategy.StrategyType, &strategy.Script)
	if err != nil {
		return nil, err
	}

	return strategy, nil
}
