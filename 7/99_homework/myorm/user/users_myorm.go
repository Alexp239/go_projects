package user

// Generated code!

import "database/sql"

var DB *sql.DB

func SetDB(db *sql.DB) {
	DB = db
}

func (data *User) FindByPK(pk uint) (err error) {
	row := DB.QueryRow("SELECT id, username, info, balance, status FROM users WHERE id= ?", pk)
	err = row.Scan(&data.ID, &data.Login, &data.Info, &data.Balance, &data.Status)
	return err
}

func (data *User) Update() (err error) {
	_, err = DB.Exec(
		"UPDATE users SET username = ?, info = ?, balance = ?, status = ? WHERE id = ?",
		data.Login, data.Info, data.Balance, data.Status, data.ID,
	)
	return err
}

func (data *User) Create() (err error) {
	result, err := DB.Exec(
		"INSERT INTO users(`username`, `info`, `balance`, `status`) VALUES (?, ?, ?, ?)",
		data.Login, data.Info, data.Balance, data.Status,
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
