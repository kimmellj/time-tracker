package main

import (
	"fmt"
	"github.com/ziutek/mymysql/mysql" // https://github.com/ziutek/mymysql
	_ "github.com/ziutek/mymysql/thrsafe" // https://github.com/ziutek/mymysql
	"os"	
)

// Contants for the database Connection
const (
	DB_ADDR  = "127.0.0.1:3306"
	DB_NAME  = "time_tracker"
	DB_USER  = "root"
	DB_PASS  = "password"
	DB_PROTO = "tcp"
)


// This function will return a connection to the database
// and exit if there are any errors
func OpenDB() mysql.Conn {
	db := mysql.New(DB_PROTO, "", DB_ADDR, DB_USER, DB_PASS, DB_NAME)

	checkError(db.Connect())

	return db
}

// This function is used to check for errors
func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Check results and exit with an error if there are any
func checkedResult(rows []mysql.Row, res mysql.Result, err error) ([]mysql.Row,
	mysql.Result) {
	checkError(err)
	return rows, res
}