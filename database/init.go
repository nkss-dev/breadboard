package database

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Init(pool *pgxpool.Pool) {
	filenames := []string{"announcement", "student", "faculty", "guild", "club", "course"}
	script := []string{}
	for _, filename := range filenames {
		// Read PostgreSQL script
		path := filepath.Join("database", "schemas", filename+".sql")
		file, ioErr := os.ReadFile(path)
		if ioErr != nil {
			log.Fatalln(ioErr)
		}

		// Execute PostgreSQL script
		script = append(script, string(file))
	}

	_, err := pool.Exec(context.Background(), strings.Join(script, "\n"))
	if err != nil {
		log.Fatalln(err)
	}
}
