package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	// This is the configured port for docker-compose
	if port == 5433 {
		ssl += `disable`
	}
	// Port for production
	if port == 5432 {
		ssl += `require`
	}
	// This was used to run docker independently.
	if port == 2222 {
		ssl += `disable`
	}
	psql_info := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=%s",
		host,
		port,
		username,
		password,
		name,
		ssl,
	)
	fmt.Println(psql_info)
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
	stmts := strings.Split(string(schemaSQL), ";")

	for i, stmt := range stmts {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		_, err := db.Exec(stmt)
		if err != nil {
			return fmt.Errorf("error executing statement %d: %v\nSQL: %s", i+1, err, stmt)
		}
	}
	// Execute permissions SQL page
	permissionPath := filepath.Join("database", "db_init_permissions.sql")
	if _, err := os.Stat(permissionPath); os.IsNotExist(err) {
		return fmt.Errorf("permissions init sql file not found %s", err)
	}
	permissionSQL, err := os.ReadFile(permissionPath)
	if err != nil {
		return fmt.Errorf("error reading init permissions file: %v", err)
	}
	_, err = db.Exec(string(permissionSQL))
	if err != nil {
		return fmt.Errorf("error executing schema SQL: %v", err)
	}

	fmt.Println("Database schema created successfully.")
	return nil
}
