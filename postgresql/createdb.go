package postgresql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/emculber/database_access/postgresql"
)

func ConnectToDatabase(dbname string, host string, port int, username string, password string) *sql.DB {
	db_url := fmt.Sprintf("postgres://%s:%s@%s/%s", username, password, host, dbname)
	fmt.Println("Connecting to database:", db_url)
	db, err := sql.Open("postgres", db_url)
	if err != nil {
		fmt.Println("Failed to connecting to database")
		fmt.Println(err)
		return nil
	}
	fmt.Println("Successfuly connected to database")
	return db
}

func CheckIfDatabaseExists(db *sql.DB, database_name string) (bool, error) {
	statement := fmt.Sprintf("select count(*) from pg_database where datname='%s'", database_name)
	fmt.Println("Checking if database exists with the statement:", statement)
	count, _, _ := postgresql_access.QueryDatabase(db, statement)
	fmt.Println("database statement count returned:", count)
	if count[0][0] == "1" {
		fmt.Println("Database exists")
		return true, nil
	}
	fmt.Println("Database does not exists")
	return false, nil
}

func CreateDatabase(db *sql.DB, database_name string) {
	statement := fmt.Sprintf("CREATE DATABASE %s", database_name)
	fmt.Println("Creating database:", statement)
	//err := postgresql_access.CreateDatabase(db, statement)
	_, err := db.Exec(statement)
	if err != nil {
		fmt.Println("Failed to creating database")
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully created database")
}

func CheckIfRoleExists(db *sql.DB, user_name string) (bool, error) {
	statement := fmt.Sprintf("select count(*) from pg_user where usename='%s'", user_name)
	fmt.Println("Checking if user exists with the statement:", statement)
	count, _, _ := postgresql_access.QueryDatabase(db, statement)
	fmt.Println("database statement count returned:", count)
	if count[0][0] == "1" {
		fmt.Println("User exists")
		return true, nil
	}
	fmt.Println("User does not exists")
	return false, nil
}

func CreateUser(db *sql.DB, user_name, user_password, user_roles string) {
	statement := fmt.Sprintf("CREATE ROLE %s WITH PASSWORD '%s' %s", user_name, user_password, user_roles)
	fmt.Println("Creating database:", statement)
	_, err := db.Exec(statement)
	if err != nil {
		fmt.Println("Failed to creating user")
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully created user")
}

func CheckIfTableExists(db *sql.DB, table_name string) (bool, error) {
	statement := fmt.Sprintf("select count(*) from pg_tables where tablename='%s'", strings.ToLower(table_name))
	fmt.Println("Checking if table exists with the statement:", statement)
	count, _, _ := postgresql_access.QueryDatabase(db, statement)
	fmt.Println("database statement count returned:", count[0][0])
	if count[0][0] == "1" {
		fmt.Println("table exists")
		return true, nil
	}
	fmt.Println("table does not exists")
	return false, nil
}
