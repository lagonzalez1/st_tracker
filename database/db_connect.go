package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"tracker/app/models"

	_ "github.com/lib/pq"
)

/*
	Getting started:
		Postgres 14 above
		Create a user with sufficient permissions to allow schema and table creation

	Downloads:
		go get "github.com/lib/pq" from root dir


		db.SetMaxOpenConns(10)                 // hard cap
		db.SetMaxIdleConns(5)                  // don’t keep too many idle
		db.SetConnMaxIdleTime(5 * time.Minute) // kill idle connections
		db.SetConnMaxLifetime(1 * time.Hour)   // recycle connections

*/

func ConnectDB(DB models.PostGresConfig) (*sql.DB, error) {
	var ssl string
	host, username := DB.Host, DB.Username
	password, name := DB.Password, DB.Name
	port, err := strconv.Atoi(DB.Port)
	if err != nil {
		return nil, err
	}
	switch port {
	case 5433:
		ssl += `disable`
		break
	case 2222:
		ssl += `disable`
		break
	case 5432:
		ssl += `require`
		break
	}
	psql_info := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=%s",
		host,
		port,
		username,
		password,
		name,
		ssl,
	)
	// Need to set max connection pool here ?
	db, err := sql.Open("postgres", psql_info)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	err = CreateSchemaIfNotExist(db)
	if err != nil {
		fmt.Println("[POSTGRES ERROR]error on createSchemaIfNotExist")
		return nil, err
	}
	fmt.Printf("[POSTGRES] Postgres connected, url: %s", psql_info)
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

	// Execute permissions SQL page
	subscriptionsPath := filepath.Join("database", "db_init_subscriptions.sql")
	if _, err := os.Stat(subscriptionsPath); os.IsNotExist(err) {
		return fmt.Errorf("permissions init sql file not found %s", err)
	}
	subscriptionsSQL, err := os.ReadFile(subscriptionsPath)
	if err != nil {
		return fmt.Errorf("error reading init permissions file: %v", err)
	}
	_, err = db.Exec(string(subscriptionsSQL))
	if err != nil {
		return fmt.Errorf("error executing schema SQL: %v", err)
	}

	fmt.Println("[POSTGRES] Database schema created successfully.")
	return nil
}
