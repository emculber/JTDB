package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/emculber/database_access/postgresql"
	"github.com/lib/pq"
)

type Configuration struct {
	Database []database
}

type database struct {
	Connection connection
	Name       string
	Tables     []table
}

type connection struct {
	Host     string
	Port     int
	Username string
	Password string
	Dbname   string
}

type table struct {
	Name    string
	Columns []column
}

type column struct {
	Name        string
	Datatype    string
	Constraints []string
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

func ConnectToDatabase(connection connection) *sql.DB {
	db_url := fmt.Sprintf("postgres://%s:%s@%s/%s", connection.Username, connection.Password, connection.Host, connection.Dbname)
	db, err := sql.Open("postgres", db_url)
	if err != nil {
		fmt.Println(err)
	}
	return db
}

func CreateDatabase(database_name string) {
	cmd := exec.Command("createdb", database_name)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Printf(out.String())
}

func CreateTableSQL(table_struct table) string {
	create_table := "CREATE TABLE " + table_struct.Name + " ("
	for i, column_struct := range table_struct.Columns {
		//fmt.Println(column_struct)
		//TODO: if columns dont exits create column
		create_table += column_struct.Name + " " + strings.ToLower(column_struct.Datatype)
		for _, constrant := range column_struct.Constraints {
			create_table += " " + strings.ToLower(constrant)
		}
		if i != len(table_struct.Columns)-1 {
			create_table += ", "
		}
	}
	create_table += ");"
	return create_table
}

func JsonToDatabse() {
	configuration := testJson()

	for _, database_struct := range configuration.Database {
		//TODO: Check if DB exists   psql -lqt | cut -d \| -f 1 | grep -qw <dbname>
		CreateDatabase(database_struct.Name)
		//TODO: Connect to database
		database_struct.Connection.Dbname = database_struct.Name
		db := ConnectToDatabase(database_struct.Connection)
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
						statement := CreateTableSQL(table_struct)
						fmt.Println(statement)
						err := postgresql_access.CreateDatabaseTable(db, statement)
						fmt.Println(err)
					}
				}
			}
		}
	}
}

func main() {
	JsonToDatabse()
}
