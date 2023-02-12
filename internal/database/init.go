package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Init(db *sql.DB) {
	filenames := []string{"announcement", "student", "faculty", "guild", "club", "course"}
	script := []string{}
	for _, filename := range filenames {
		// Read PostgreSQL script
		path := filepath.Join("sql", "schemas", filename+".sql")
		file, ioErr := os.ReadFile(path)
		if ioErr != nil {
			log.Fatalln(ioErr)
		}

		// Execute PostgreSQL script
		script = append(script, string(file))
	}

	_, err := db.Exec(strings.Join(script, "\n"))
	if err != nil {
		log.Fatalln(err)
	}
}
