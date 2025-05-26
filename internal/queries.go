package internal

import "database/sql"

type DBValues struct {
	Month      string
	MonthOf    string
	Day        string
	Dedication string
}

type DBModel struct {
	DB *sql.DB
}

func (m DBModel) GetDbValues(day, month int) (*DBValues, error) {
	query := `
		SELECT m.name, m.meaning, d.name, dm.dedication
		FROM days_months dm
		INNER JOIN days d ON dm.day_id = d.id
		INNER JOIN months m ON dm.month_id = m.id
		WHERE d.id = $1 and m.id = $2
	`

	var vals DBValues

	err := m.DB.QueryRow(query, day, month).Scan(
		&vals.Month,
		&vals.MonthOf,
		&vals.Day,
		&vals.Dedication,
	)

	if err != nil {
		return nil, err
	}

	return &vals, nil
}
