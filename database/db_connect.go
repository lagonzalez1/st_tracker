package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"tracker/app/config"

	_ "github.com/lib/pq"
)

/*
	Getting started:
		Postgres 14 above
		Create a user with sufficient permissions to allow schema and table creation

	Downloads:
		go get "github.com/lib/pq" from root dir

*/

func ConnectDB() (*sql.DB, error) {
	env_config, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	host, username := env_config.DB.Host, env_config.DB.Username
	password, name := env_config.DB.Password, env_config.DB.Name
	port, err := strconv.Atoi(env_config.DB.Port)
	if err != nil {
		return nil, err
	}
	psql_info := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		host,
		port,
		username,
		password,
		name)
	db, err := sql.Open("postgres", psql_info)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Postgres connected.")
	return db, nil
}

// @Params: db pointer @Returns: Possible error
func CreateSchemaIfNotExist(db *sql.DB) error {
	schemaPath := filepath.Join("database", "db_schema.sql")
	schemaSQL, err := os.ReadFile(schemaPath)

	if err != nil {
		return fmt.Errorf("error reading schema file: %v", err)
	}
	_, err = db.Exec(string(schemaSQL))
	if err != nil {
		return fmt.Errorf("error executing schema %v", err)
	}
	return nil
}
