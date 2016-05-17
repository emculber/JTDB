package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/emculber/database_access/postgresql"

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

func GetDatabaseConnection(database_name string, user User, connection Connection) *sql.DB {
	fmt.Println("Parsing Database connection information")
	dbname := strings.Split(database_name, "{")
	dbuser := strings.Split(user.Username, "{")

	fmt.Println("Connecting to default database:",
		"\nDatabase Name=", dbname[0],
		"\nDatabase Host=", connection.Host,
		"\nDatabase Port=", connection.Port,
		"\nDatabase User=", dbuser[0],
		"\nDatabase Password=", user.Password)

	db := postgresql.ConnectToDatabase(
		dbname[0],
		connection.Host,
		connection.Port,
		dbuser[0],
		user.Password,
	)
	return db
}

func CreateTableSQL(table_struct Table) string {
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

func JsonToDatabse(configuration Configuration) Configuration {
	fmt.Println("Traversing database structure")
	for i, database_struct := range configuration.Database {
		fmt.Println("Getting database connection")
		var db *sql.DB
		fmt.Println("Checking if Default Database is set:", database_struct.Default)
		if database_struct.Default != "" {
			db = GetDatabaseConnection(database_struct.Default, database_struct.Users.DefaultUser, database_struct.Connection)
		} else {
			db = GetDatabaseConnection(database_struct.Name, database_struct.Users.DefaultUser, database_struct.Connection)
		}

		fmt.Println("Checking if Roles exist")
		for u_index, user := range database_struct.Users.User {
			user_name := strings.Split(user.Username, "{")
			user_roles := strings.Replace(user.Role, " ", "", -1)
			user_roles = strings.Replace(user_roles, ",", " ", -1)

			fmt.Println("Checking if user exitst:", user_name[0])
			user_exist, _ := postgresql.CheckIfRoleExists(db, user_name[0])
			//TODO: Check if roles match
			if !user_exist {
				fmt.Println("User does not exits, Creating User")
				postgresql.CreateUser(db, user_name[0], user.Password, user_roles)
			}
			if strings.Contains(user_name[1], "Default") {
				fmt.Println("Starting Proccess of setting user to default")
				defaut_user_name := strings.Split(database_struct.Users.DefaultUser.Username, "{")
				if strings.Contains(defaut_user_name[1], "Remove") {
					fmt.Println("Current default user is set to be removed")
					configuration.Database[i].Users.DefaultUser.Username = user_name[0]
					configuration.Database[i].Users.DefaultUser.Password = user.Password
					configuration.Database[i].Users.DefaultUser.Role = user.Role
					fmt.Println("New default",
						"\nusername:", configuration.Database[i].Users.DefaultUser.Username,
						"\nPassword:", configuration.Database[i].Users.DefaultUser.Password,
						"\nRole:", configuration.Database[i].Users.DefaultUser.Role)
					fmt.Println("Deleting Old User")
					configuration.Database[i].Users.User = append(configuration.Database[i].Users.User[:u_index], configuration.Database[i].Users.User[u_index+1:]...)
				} else {
					fmt.Println("default user is not set to be removed no action will be taken")
				}
			}
		}

		if database_struct.Default != "" {
			fmt.Println("Closing database connection for new connection")
			db.Close()
			fmt.Println("Openning new connection")
			db = GetDatabaseConnection(configuration.Database[i].Default, configuration.Database[i].Users.DefaultUser, configuration.Database[i].Connection)

			fmt.Println("Checking if Database:", database_struct.Name, "Exists")
			exist, _ := postgresql.CheckIfDatabaseExists(db, database_struct.Name)

			if !exist {
				fmt.Println("Database does not exits, Creating database")
				postgresql.CreateDatabase(db, database_struct.Name)
			}
		}
		fmt.Println("Closing database connection for new connection")
		db.Close()
		fmt.Println("Openning new connection")
		db = GetDatabaseConnection(configuration.Database[i].Name, configuration.Database[i].Users.DefaultUser, configuration.Database[i].Connection)
		fmt.Println("Checking and creating tables")
		for _, table_struct := range database_struct.Tables {
			table_exist, _ := postgresql.CheckIfTableExists(db, table_struct.Name)
			if !table_exist {
				fmt.Println("Table does not exits, Creating User")
				statement := CreateTableSQL(table_struct)
				fmt.Println("Created statment to create table:", statement)
				err := postgresql_access.CreateDatabaseTable(db, statement)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}

	return configuration
}

func main() {

	path := "/home/erik/programming/golang/src/Scheduler"

	fmt.Println("Running Json to Database")

	fmt.Println("Loading Json Configuration")
	configuration := testJson(path)
	new_configuration := JsonToDatabse(testJson(path))

	if !reflect.DeepEqual(configuration, new_configuration) {
		fmt.Println("Configurations are not equal")
		err := os.Mkdir(path+"/PastConfig", 0711)
		if err != nil {
			fmt.Println(err)
		}
		files, _ := ioutil.ReadDir(path + "/PastConfig")
		f, err := os.Create(path + "/PastConfig/db.conf_v" + strconv.Itoa(len(files)) + ".json")
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()

		configuration_json, err := json.Marshal(configuration)
		if err != nil {
			fmt.Println(err)
		}
		f.Write(configuration_json)

		err = os.Remove(path + "/db.conf.json")
		if err != nil {
			fmt.Println(err)
		}

		conf, err := os.Create(path + "/db.conf.json")
		if err != nil {
			fmt.Println(err)
		}
		defer conf.Close()
		new_configuration_json, err := json.Marshal(new_configuration)
		if err != nil {
			fmt.Println(err)
		}
		conf.Write(new_configuration_json)
	} else {
		fmt.Println("Configurations are equal")
	}
}

func testJson(path string) Configuration {
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
