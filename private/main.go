package main

import (
	"code.google.com/p/goweb/goweb"
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/thrsafe"
	"os"
)

const (
	DB_ADDR  = "127.0.0.1:3306"
	DB_NAME  = "time_tracker"
	DB_USER  = "root"
	DB_PASS  = "fawkes"
	DB_PROTO = "tcp"
)

func OpenDB() mysql.Conn {
	db := mysql.New(DB_PROTO, "", DB_ADDR, DB_USER, DB_PASS, DB_NAME)

	fmt.Printf("Connect to %s:%s... ", DB_PROTO, DB_ADDR)
	checkError(db.Connect())
	printOK()

	return db
}

func printOK() {
	fmt.Println("OK")
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func checkedResult(rows []mysql.Row, res mysql.Result, err error) ([]mysql.Row,
	mysql.Result) {
	checkError(err)
	return rows, res
}

type Client struct {
	ID   int
	Name string
}

type ClientsController struct{}

func (cr *ClientsController) Create(cx *goweb.Context) {
	fmt.Fprintf(cx.ResponseWriter, "Create new resource")
}
func (cr *ClientsController) Delete(id string, cx *goweb.Context) {
	fmt.Fprintf(cx.ResponseWriter, "Delete resource %s", id)
}
func (cr *ClientsController) Read(id string, cx *goweb.Context) {
	fmt.Fprintf(cx.ResponseWriter, "Read resource %s", id)
}
func (cr *ClientsController) ReadMany(cx *goweb.Context) {
	db := OpenDB()

	rows, res := checkedResult(db.Query("select * from clients"))

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
func (cr *ClientsController) Update(id string, cx *goweb.Context) {

}

func main() {
	// create our resource controller
	controller := new(ClientsController)

	// make the mapping
	goweb.MapRest("/clients", controller)

    goweb.ConfigureDefaultFormatters()
    goweb.ListenAndServe(":8080")


}