package init_db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func InitializeDB() *sql.DB {

	dbDir := filepath.Join("pkg", "db")
	dbFile := filepath.Join(dbDir, "scheduler.db")

	if _, err := os.Stat(dbDir); os.IsNotExist(err) {

		err = os.Mkdir(dbDir, 0755)
		if err != nil {
			log.Fatal("error creating database directory: ", err)
		}
	}

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		file, err := os.Create(dbFile)
		if err != nil {
			log.Fatal("error creating database file: ", err)
		}
		file.Close()

		err = createTable(dbFile)
		if err != nil {
			log.Fatal("error initializing database table: ", err)
		}

		fmt.Println("Database and table created")
	} else if err != nil {
		log.Fatal("error checking db file: ", err)
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal("error opening db: ", err)
	}

	return db
}

func createTable(dbFile string) error {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	createTableSQl := `
	CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL,
    title TEXT NOT NULL,
    comment TEXT,
    repeat TEXT
);`

	_, err = db.Exec(createTableSQl)
	return err
}
