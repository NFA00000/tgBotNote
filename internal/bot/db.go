package bot

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./notes.db")
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных: ", err)
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		text TEXT NOT NULL,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = DB.Exec(sqlStmt)
	if err != nil {
		log.Fatal("Ошибка создания таблицы: ", err)
	}
}

func AddNote(userID int64, title, text string) {
	_, err := DB.Exec("INSERT INTO notes (user_id, title, text) VALUES (?, ?, ?)", userID, title, text)
	if err != nil {
		log.Fatal("Ошибка добавления заметки: ", err)
	}
}
