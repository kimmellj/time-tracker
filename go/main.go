/*
This package creates a RESTful API for a time tracking system
*/

package main

import (
	"code.google.com/p/goweb/goweb" // http://code.google.com/p/goweb/
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

type Client struct {
	ID   int
	Name string
}

// Defining the the Clients Controller so that it can be extended
type CRUDController struct{}

// Create a client this responds to a POST request on /clients
func (cr *CRUDController) Create(cx *goweb.Context) {
	db := OpenDB()

	cx.Request.ParseForm()

	stmt, err := db.Prepare("INSERT INTO `clients` SET `name`=?")

	if err != nil {
		panic(err)
	}

	_, err = stmt.Run(cx.Request.Form["name"][0])

	if err != nil {
		cx.RespondWithErrorMessage("Trouble adding the client.", 200)
	} else {
		cx.RespondWithOK()
	}
}

// Delete a specified client this responds to a DELETE request on /clients/{id}
func (cr *CRUDController) Delete(id string, cx *goweb.Context) {
	db := OpenDB()

	stmt, err := db.Prepare("DELETE FROM `clients` WHERE `id`=?")

	if err != nil {
		cx.RespondWithErrorMessage("Trouble preparing to delete the client.", 200)
		panic(err)
	}

	_, err = stmt.Run(id)

	if err != nil {
		cx.RespondWithErrorMessage("Trouble deleting the client.", 200)
		panic(err)
	} else {
		cx.RespondWithOK()
	}
}

// Get a specifed clients details this responds to a GET request on /clients/{id}
func (cr *CRUDController) Read(queryID string, cx *goweb.Context) {
	db := OpenDB()

	stmt, err := db.Prepare("SELECT `id`, `name` FROM `clients` WHERE `id`=?")

	if err != nil {
		cx.RespondWithErrorMessage("Trouble preparing to fetch the client.", 200)
		panic(err)
	}

	result, err := stmt.Run(queryID)

	if err != nil {
		cx.RespondWithErrorMessage("Trouble fetching the client.", 200)
		panic(err)
	} else {
		id := result.Map("id")
		name := result.Map("name")

		row, err := result.GetRow()

		if err != nil {
			panic(err)
		}

		cx.RespondWithData(Client{row.Int(id), row.Str(name)})
	}
}

// Get all of the clients this responds to a GET request on /clients
func (cr *CRUDController) ReadMany(cx *goweb.Context) {
	db := OpenDB()

	rows, res := checkedResult(db.Query("SELECT id,name FROM clients"))

	id := res.Map("id")
	name := res.Map("name")

	numRows := len(rows)

	var clients []Client
	clients = make([]Client, numRows, numRows)

	for ii, row := range rows {
		clients[ii] = Client{row.Int(id), row.Str(name)}
	}

	cx.RespondWithData(clients)
}

// Update a specified client this responds to a PUT on /clients/{id} with form data
func (cr *CRUDController) Update(id string, cx *goweb.Context) {
	db := OpenDB()

	cx.Request.ParseForm()

	stmt, err := db.Prepare("UPDATE `clients` SET `name`=? WHERE `id`=?")

	if err != nil {
		panic(err)
	}

	_, err = stmt.Run(cx.Request.Form["name"][0], id)

	if err != nil {
		cx.RespondWithErrorMessage("Trouble updating the client.", 200)
	} else {
		cx.RespondWithOK()
	}
}

// Main Application loop, listens for web requests and then forwards them to goweb
func main() {
	// create our resource controller
	clientsController := new(CRUDController)
	projectsController := new(CRUDController)

	// make the mapping
	goweb.MapRest("/clients", clientsController)
	goweb.MapRest("/projects", projectsController)

	goweb.ConfigureDefaultFormatters()
	goweb.ListenAndServe(":8080")

}
