package main

import (
	"fmt"
	"strconv"
	"net/http"
	"github.com/ant0ine/go-json-rest"
)

type Book struct {
	Id int
	Book_id string
	Book_title string
	Book_author string
	Book_category_num string
	Book_isbn string
	Book_page_num string
	Book_price string
	Book_pubdate string
	Book_publisher string
	Book_ref_no string
	Book_summary string
}

type Status struct {
	Status_code int
	Status_desc string
}

type GetRespone struct {
	Book Book
	Status Status
}

func main() {
	parseConfig("./conf/config.conf")
	printVersion()
	fmt.Println("Start atat server......")
	printAsciiIcon()
//	count, err := dbGetBookCount()
	handler := rest.ResourceHandler{}
	handler.SetRoutes(
		rest.Route{"GET", "/book/:id", GetBook},
	)
	http.ListenAndServe(":8080", &handler)
}


func GetBook(w *rest.ResponseWriter, req *rest.Request) {
//	params := map[string][]string{"userId":[]string{ req.PathParam("id") }}
	id, err := strconv.Atoi(req.PathParam("id"))
	checkErr(err)
	book, err := dbGetBookFromBookId(id)
	checkErr(err)
	if book.Id != 0 {
		w.WriteJson(&book)
	}else {
		w.WriteHeader(404)
	}
}


















