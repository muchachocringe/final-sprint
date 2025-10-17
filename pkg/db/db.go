package db

import (
	//"database/sql"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var db *sqlx.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL,
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(255) NOT NULL DEFAULT ""
);

CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	var database *sqlx.DB
	database, err = sqlx.Connect("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД: %v", err)
	}

	if install {
		if _, err := database.Exec(schema); err != nil {
			return fmt.Errorf("ошибка создания таблиц: %v", err)
		}
	}

	db = database
	return nil
}

func GetDB() *sqlx.DB {
	return db
}
