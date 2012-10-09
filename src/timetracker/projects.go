/*
/*

Example URLS:

GET: http://localhost:8080/projects.json

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

GET: http://localhost:8080/projects/1.json

{
    "C": "",
    "S": 200,
    "D": {
        "ID": 1,
        "Name": "Jansoort"
    },
    "E": null
}

POST: http://localhost:8080/projects.json

Will create a project with the POST data sent

DELETE: http://localhost:8080/projects/1.json

Will delete the project with id 1

PUT: GET: http://localhost:8080/projects/1.json

Will update the project with id 1 with the POST data sent
*/


package main

import (
	"code.google.com/p/goweb/goweb" // http://code.google.com/p/goweb/
)

// This struct represents a project
type Project struct {
	ID   int
	Name string
}

// Defining the the Projects Controller so that it can be extended
type ProjectsController struct{}

// Create a project this responds to a POST request on /projects
func (cr *ProjectsController) Create(cx *goweb.Context) {
	db := OpenDB()

	cx.Request.ParseForm()

	stmt, err := db.Prepare("INSERT INTO `projects` SET `name`=?")

	if err != nil {
		panic(err)
	}

	_, err = stmt.Run(cx.Request.Form["name"][0])

	if err != nil {
		cx.RespondWithErrorMessage("Trouble adding the project.", 200)
	} else {
		cx.RespondWithOK()
	}
}

// Delete a specified project this responds to a DELETE request on /projects/{id}
func (cr *ProjectsController) Delete(id string, cx *goweb.Context) {
	db := OpenDB()

	stmt, err := db.Prepare("DELETE FROM `projects` WHERE `id`=?")

	if err != nil {
		cx.RespondWithErrorMessage("Trouble preparing to delete the project.", 200)
		panic(err)
	}

	_, err = stmt.Run(id)

	if err != nil {
		cx.RespondWithErrorMessage("Trouble deleting the project.", 200)
		panic(err)
	} else {
		cx.RespondWithOK()
	}
}

// Get a specifed projects details this responds to a GET request on /projects/{id}
func (cr *ProjectsController) Read(queryID string, cx *goweb.Context) {
	db := OpenDB()

	stmt, err := db.Prepare("SELECT `id`, `name` FROM `projects` WHERE `id`=?")

	if err != nil {
		cx.RespondWithErrorMessage("Trouble preparing to fetch the project.", 200)
		panic(err)
	}

	result, err := stmt.Run(queryID)

	if err != nil {
		cx.RespondWithErrorMessage("Trouble fetching the project.", 200)
		panic(err)
	} else {
		id := result.Map("id")
		name := result.Map("name")

		row, err := result.GetRow()

		if err != nil {
			panic(err)
		}

		cx.RespondWithData(Project{row.Int(id), row.Str(name)})
	}
}

// Get all of the projects this responds to a GET request on /projects
func (cr *ProjectsController) ReadMany(cx *goweb.Context) {
	db := OpenDB()

	rows, res := checkedResult(db.Query("SELECT id,name FROM projects"))

	id := res.Map("id")
	name := res.Map("name")

	numRows := len(rows)

	var projects []Project
	projects = make([]Project, numRows, numRows)

	for ii, row := range rows {
		projects[ii] = Project{row.Int(id), row.Str(name)}
	}

	cx.RespondWithData(projects)
}

// Update a specified project this responds to a PUT on /projects/{id} with form data
func (cr *ProjectsController) Update(id string, cx *goweb.Context) {
	db := OpenDB()

	cx.Request.ParseForm()

	stmt, err := db.Prepare("UPDATE `projects` SET `name`=? WHERE `id`=?")

	if err != nil {
		panic(err)
	}

	_, err = stmt.Run(cx.Request.Form["name"][0], id)

	if err != nil {
		cx.RespondWithErrorMessage("Trouble updating the project.", 200)
	} else {
		cx.RespondWithOK()
	}
}