package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	common "github.com/n7down/ssh-chess/common"
)

var (
	db *sql.DB
)

func initDB() (err error) {
	dbUser := "root"
	dbPassword := "password"
	dbHost := common.GetEnv("DB_HOST", "localhost")
	dbPort := common.GetEnv("DB_PORT", "5432")
	dbName := "chess"
	connectionString := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}
	fmt.Println("connected to: " + connectionString)
	return nil
}

func GetDB() *sql.DB {
	if db == nil {
		if err := initDB(); err != nil {
			fmt.Printf(fmt.Sprintf("Error with the database: %s", err.Error()))
		}
	}
	return db
}
