package database

import (
	"database/sql"
	"fmt"
	"log"
)

type File struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Data      []byte `json:"data"`
	Extension string `json:"extension"`
}

func InitDB(db *sql.DB) {

	create_query := `CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		extension TEXT NOT NULL,
		data BLOB,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		modified_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	revert_query := `DROP TABLE IF EXISTS files;`

	trigger_query := `CREATE TRIGGER IF NOT EXISTS update_modified_at
	AFTER UPDATE ON files
	FOR EACH ROW
	BEGIN
		UPDATE files SET modified_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
	END;`

	revert_trigger_query := `DROP TRIGGER IF EXISTS update_modified_at;`

	fmt.Println("Creating the database")
	_, err := db.Exec(create_query)
	if err != nil {
		db.Exec(revert_query)
		log.Fatal("Error creating table files ", err, " all creations have been reverted")
	}

	_, err = db.Exec(trigger_query)

	if err != nil {
		db.Exec(revert_trigger_query)
		db.Exec(revert_query)
		log.Fatal("Error creating trigger query ", err, " all creations have been reverted")
	}
}
