package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"../JTDB/postgresql"
)

type Configuration struct {
	Database []database
}

type database struct {
	Default    string
	Name       string
	Connection Connection
	Users      Users
	Tables     []Table
}

type Connection struct {
	Host string
	Port int
}

type Users struct {
	DefaultUser User `json:"Default User"`
	User        []User
}

type User struct {
	Username string
	Password string
	Role     string
}

type Table struct {
	Name    string
	Columns []Column
}

type Column struct {
	Name        string
	Datatype    string
	Constraints []string
}

func GetDatabaseConnection(database_struct database) *sql.DB {
	fmt.Println("Parsing Database connection information")
	//database_name := strings.Split(database_struct.Name, "{")
	database_user := strings.Split(database_struct.Users.DefaultUser.Username, "{")
	database_password := strings.Split(database_struct.Users.DefaultUser.Password, "{")

	fmt.Println("Connecting to default database:",
		"\nDatabase Name=", database_struct.Default,
		"\nDatabase Host=", database_struct.Connection.Host,
		"\nDatabase Port=", database_struct.Connection.Port,
		"\nDatabase User=", database_user[0],
		"\nDatabase Password=", database_password[0])

	db := postgresql.ConnectToDatabase(
		database_struct.Default,
		database_struct.Connection.Host,
		database_struct.Connection.Port,
		database_user[0],
		database_password[0],
	)
	return db
}

func JsonToDatabse() {
	fmt.Println("Loading Json Configuration")
	configuration := testJson()

	fmt.Println("Traversing database structure")
	for _, database_struct := range configuration.Database {
		fmt.Println("Getting database connection")
		db := GetDatabaseConnection(database_struct)

		fmt.Println("Checking if Roles exist")
		for _, user := range database_struct.Users.User {
			user_name := strings.Split(user.Username, "{")
			user_roles := strings.Replace(user.Role, " ", "", -1)
			user_roles = strings.Replace(user_roles, ",", " ", -1)

			fmt.Println("Checking if user exitst:", user_name[0])
			user_exist, _ := postgresql.CheckIfRoleExists(db, user_name[0])
			//TODO: Check if roles match
			if !user_exist {
				fmt.Println("User does not exits, Creating User")
				postgresql.CreateUser(db, user_name[0], user_roles)
			}
		}
		return

		fmt.Println("Checking if Database:", database_struct.Name, "Exists")
		exist, _ := postgresql.CheckIfDatabaseExists(db, database_struct.Name)

		if !exist {
			fmt.Println("Database does not exits, Creating database")
			postgresql.CreateDatabase(db, database_struct.Name)
		}
	}

	return
	/*
		for _, database_struct := range configuration.Database {
			//TODO: Check if DB exists   psql -lqt | cut -d \| -f 1 | grep -qw <dbname>
			postgresql.CreateDatabase(database_struct.Name)
			//TODO: Connect to database
			database_struct.Connection.Dbname = database_struct.Name
			db := postgresql.ConnectToDatabase(database_struct.Connection)
			//fmt.Println(db)
			for _, table_struct := range database_struct.Tables {
				//TODO: Check if tables exist
				statement := "Select * from " + table_struct.Name
				err := postgresql_access.CreateDatabaseTable(db, statement)
				if err != nil {
					if err, ok := err.(*pq.Error); ok {
						//fmt.Println(err.Code)
						if err.Code == "42P01" {
							//fmt.Println(err.Code)
							//TODO: if Table does not exist create Table
							statement := postgresql.CreateTableSQL(table_struct)
							fmt.Println(statement)
							err := postgresql_access.CreateDatabaseTable(db, statement)
							fmt.Println(err)
						}
					}
				}
			}
		}
	*/
}

func main() {
	fmt.Println("Running Json to Database")
	JsonToDatabse()
}

func testJson() Configuration {
	path, _ := os.Getwd()
	fmt.Println("Loading test Json Configuration at path: ", path)

	config_file, err := os.Open(path + "/db.conf.json")
	if err != nil {
		fmt.Println("Test Json Configuration failed to loaded successfully")
		fmt.Println(err)
	}
	fmt.Println("Test Json Configuration was loaded successfully")

	var configuration Configuration

	fmt.Println("Decoding json")
	decoder := json.NewDecoder(config_file)
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("Failed to Decoding json")
		fmt.Println(err)
	}
	fmt.Println("Decoding json was sucessful")
	//fmt.Println(configuration)
	return configuration
}
