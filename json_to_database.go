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
	database_name := strings.Split(database_struct.Name, "{")
	database_user := strings.Split(database_struct.Users.DefaultUser.Username, "{")
	database_password := strings.Split(database_struct.Users.DefaultUser.Password, "{")

	db := postgresql.ConnectToDatabase(
		database_name[0],
		database_struct.Connection.Host,
		database_struct.Connection.Port,
		database_user[0],
		database_password[0],
	)
	return db
}

func JsonToDatabse() {
	configuration := testJson()

	for _, database_struct := range configuration.Database {
		db := GetDatabaseConnection(database_struct)

		exist, _ := postgresql.CheckIfDatabaseExists(db, database_struct.Name)

		if !exist {
			fmt.Println(exist)
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
	JsonToDatabse()
}

func testJson() Configuration {
	path, _ := os.Getwd()

	//dat, _ := ioutil.ReadFile(path + "/db.conf.json")
	//fmt.Print(string(dat))

	config_file, err := os.Open(path + "/db.conf.json")
	if err != nil {
		fmt.Println(err)
	}

	var configuration Configuration

	decoder := json.NewDecoder(config_file)
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(configuration)
	return configuration
}
