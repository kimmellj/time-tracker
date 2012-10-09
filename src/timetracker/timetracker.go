/*
This package creates a RESTful API for a time tracking system
*/

package main

import (
	"code.google.com/p/goweb/goweb" // http://code.google.com/p/goweb/
)

// Main Application loop, listens for web requests and then forwards them to goweb
func main() {
	// create our resource controller
	clientsController := new(ClientsController)
	projectsController := new(ProjectsController)

	// make the mapping
	goweb.MapRest("/clients", clientsController)
	goweb.MapRest("/projects", projectsController)

	goweb.ConfigureDefaultFormatters()
	goweb.ListenAndServe(":8080")

}
