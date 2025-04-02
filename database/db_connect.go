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
	var ssl string
	host, username := env_config.DB.Host, env_config.DB.Username
	password, name := env_config.DB.Password, env_config.DB.Name
	port, err := strconv.Atoi(env_config.DB.Port)
	if err != nil {
		return nil, err
	}
	if port == 5433 {
		ssl += `disable`
	}
	if port == 5432 {
		ssl += `require`
	}
	psql_info := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=%s",
		host,
		port,
		username,
		password,
		name,
		ssl,
	)
	db, err := sql.Open("postgres", psql_info)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	err = CreateSchemaIfNotExist(db)
	if err != nil {
		fmt.Println("error on createSchemaIfNotExist")
		return nil, err
	}
	fmt.Println("Postgres connected.")
	return db, nil
}

// @Params: db pointer @Returns: Possible error
func CreateSchemaIfNotExist(db *sql.DB) error {
	var tableExists bool
	// Check if table exists
	err := db.QueryRow(`
        SELECT EXISTS (
            SELECT 1 FROM information_schema.tables 
            WHERE table_schema = 'stu_tracker' AND table_name = 'organization'
        );`).Scan(&tableExists)

	if err != nil {
		return fmt.Errorf("error checking table existence: %v", err)
	}

	// Skip table creation if it already exists
	if tableExists {
		fmt.Println("Table 'stu_tracker' already exists.")
		return nil
	}

	// Debug: Verify file exists before reading
	schemaPath := filepath.Join("database", "db_schema.sql")
	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		return fmt.Errorf("schema file not found: %s", schemaPath)
	}

	// Read SQL file
	schemaSQL, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("error reading schema file: %v", err)
	}

	// Execute SQL script
	_, err = db.Exec(string(schemaSQL))
	if err != nil {
		return fmt.Errorf("error executing schema SQL: %v", err)
	}

	fmt.Println("Database schema created successfully.")
	return nil
}
