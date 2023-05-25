package utils

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"
)

func InitDb() *sql.DB {
	err := godotenv.Load(".env")
	CheckError(err)

	var (
		host = os.Getenv("PSQL_HOST")
		// port     = os.Getenv("PSQL_PORT")
		user     = os.Getenv("PSQL_USER")
		password = os.Getenv("PSQL_PWD")
		dbname   = "arthurdb"
	)

	connStr := "postgres://" + user + ":" + password + "@" + host + "/" + dbname + "?sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	CheckError(err)

	return db
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
