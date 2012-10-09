/*
/*

Example URLS:

GET: http://localhost:8080/clients.json

{
    "C": "",
    "S": 200,
    "D": [
        {
            "ID": 1,
            "Name": "Jansoort"
        },
        {
            "ID": 2,
            "Name": "Reebok"
        },
        {
            "ID": 3,
            "Name": "Tasha"
        }
    ],
    "E": null
}

GET: http://localhost:8080/clients/1.json

{
    "C": "",
    "S": 200,
    "D": {
        "ID": 1,
        "Name": "Jansoort"
    },
    "E": null
}

POST: http://localhost:8080/clients.json

Will create a client with the POST data sent

DELETE: http://localhost:8080/clients/1.json

Will delete the client with id 1

PUT: GET: http://localhost:8080/clients/1.json

Will update the client with id 1 with the POST data sent
*/


package main

import (
	"code.google.com/p/goweb/goweb" // http://code.google.com/p/goweb/
)


// This struct represents a client
type Client struct {
	ID   int
	Name string
}

// Defining the the Clients Controller so that it can be extended
type ClientsController struct{}

// Create a client this responds to a POST request on /clients
func (cr *ClientsController) Create(cx *goweb.Context) {
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
func (cr *ClientsController) Delete(id string, cx *goweb.Context) {
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
func (cr *ClientsController) Read(queryID string, cx *goweb.Context) {
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
func (cr *ClientsController) ReadMany(cx *goweb.Context) {
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
func (cr *ClientsController) Update(id string, cx *goweb.Context) {
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