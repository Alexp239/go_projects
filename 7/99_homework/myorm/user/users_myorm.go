package user

// Generated code!

import "database/sql"

var DB *sql.DB

func SetDB(db *sql.DB) {
	DB = db
}

func (data *User) FindByPK(pk uint) (err error) {
	var info_null *string
	row := DB.QueryRow("SELECT id, username, info, balance, status FROM users WHERE id= ?", pk)
	err = row.Scan(&data.ID, &data.Login, &info_null, &data.Balance, &data.Status)
	if info_null != nil {
		data.Info = *info_null
	} else {
		data.Info = ""
	}
	return err
}

func (data *User) Update() (err error) {
	var info_null sql.NullString
	if len(data.Info) == 0 {
		info_null = sql.NullString{}
	} else {
		info_null = sql.NullString{
			String: data.Info,
			Valid:  true,
		}
	}
	_, err = DB.Exec(
		"UPDATE users SET username = ?, info = ?, balance = ?, status = ? WHERE id = ?",
		data.Login, info_null, data.Balance, data.Status, data.ID,
	)
	return err
}

func (data *User) Create() (err error) {
	var info_null sql.NullString
	if len(data.Info) == 0 {
		info_null = sql.NullString{}
	} else {
		info_null = sql.NullString{
			String: data.Info,
			Valid:  true,
		}
	}
	result, err := DB.Exec(
		"INSERT INTO users(`username`, `info`, `balance`, `status`) VALUES (?, ?, ?, ?)",
		data.Login, info_null, data.Balance, data.Status,
	)
	if err != nil {
		return
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return
	}
	data.ID = uint(lastID)
	return nil
}
