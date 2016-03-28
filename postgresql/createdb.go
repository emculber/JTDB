package postgresql

import (
	"database/sql"
	"fmt"

	"github.com/emculber/database_access/postgresql"
)

func ConnectToDatabase(dbname string, host string, port int, username string, password string) *sql.DB {
	db_url := fmt.Sprintf("postgres://%s:%s@%s/%s", username, password, host, dbname)
	fmt.Println(db_url)
	db, err := sql.Open("postgres", db_url)
	if err != nil {
		fmt.Println(err)
	}
	return db
}

func CheckIfDatabaseExists(db *sql.DB, database_name string) (bool, error) {
	statement := fmt.Sprintf("select count(*) from pg_database where datname='%s'", database_name)
	fmt.Println(statement)
	_, count, _ := postgresql_access.QueryDatabase(db, statement)
	if count == 1 {
		return true, nil
	}
	return false, nil
}

func CreateDatabase(db *sql.DB, database_name string) {
	statement := fmt.Sprintf("CREATE DATABASE %s", database_name)
	err := postgresql_access.CreateDatabase(db, statement)
	fmt.Println(err)
}

/*
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
}*/
