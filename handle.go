package main

import (
	"github.com/ant0ine/go-json-rest"
)

// Store single book infomation.
type Book struct {
	Id                int
	Book_id           string
	Book_title        string
	Book_author       string
	Book_category_num string
	Book_isbn         string
	Book_page_num     string
	Book_price        string
	Book_pubdate      string
	Book_publisher    string
	Book_ref_no       string
	Book_summary      string
}

// Store multiple books.
type BookList struct {
	Start int
	Count int
	Total int
	Books []Book
}

// Use for return status if some exception occurs.
type ResourceStatus struct {
	Msg     string
	Code    int
	Request string
}

// Return a ResourceStatus struct if Resource not found.
func ResourceNotFound(req *rest.Request) *ResourceStatus {
	return &ResourceStatus{
		Msg:"Resource not found",
		Code: 404,
		Request: req.Method + " " + req.URL.Path + req.URL.RawQuery,
	}
}

// Handler for route /book/:id
// Query a book_id and response one single book information in a json.
func GetBookFromBookId(w *rest.ResponseWriter, req *rest.Request) {
	//	params := map[string][]string{"userId":[]string{ req.PathParam("id") }}
	id := req.PathParam("id")
	//	query := req.URL.Query()
	//	fields := query["fields"][0]
	//	fmt.Println(query)
	if len(id) != 8 {
		var statusMsg ResourceStatus = *ResourceNotFound(req)
		w.WriteHeader(404)
		w.WriteJson(&statusMsg)
		return
	}
	//	checkErr(err)
	book, err := dbQueryBookFromBookId(id)
	checkErr(err)

	if book.Id != 0 {
		w.WriteJson(&book)
		return
	} else {
		var statusMsg ResourceStatus = *ResourceNotFound(req)
		w.WriteHeader(404)
		w.WriteJson(&statusMsg)
		return
	}
}

// Handler for route /book/isbn/:isbn
// Query from book_isbn and response one single book information in a json.
func GetBookFromBookISBN(w *rest.ResponseWriter, req *rest.Request) {
	//	params := map[string][]string{"userId":[]string{ req.PathParam("id") }}
	book_isbn := req.PathParam("isbn")
	//	query := req.URL.Query()
	//	fields := query["fields"][0]
	//	fmt.Println(query)
	//	checkErr(err)

	// Check the isbn parameter is valid
	if !isValidIsbn13(book_isbn) {
		var statusMsg ResourceStatus = *ResourceNotFound(req)
		w.WriteHeader(404)
		w.WriteJson(&statusMsg)
		return
	}
	
	book, err := dbQueryBookFromBookIsbn(book_isbn)
	checkErr(err)

	if book.Id != 0 {
		w.WriteJson(&book)
		return
	} else {
		var statusMsg ResourceStatus = *ResourceNotFound(req)
		w.WriteHeader(404)
		w.WriteJson(&statusMsg)
		return
	}
}

// Handler for route /book/search/
// Query from keyword and response a booklist with all book 
// match keyword in a json.
func GetBookListFromKeyword(w *rest.ResponseWriter, req *rest.Request) {
	parameter := req.URL.Query()
	qmap := parameter["q"]
	startmap := parameter["start"]
	countmap := parameter["count"]
	var keyword, start, count string

	if len(qmap) == 0 {
		var statusMsg ResourceStatus = *ResourceNotFound(req)
		w.WriteHeader(404)
		w.WriteJson(&statusMsg)
		return
	} else {
		keyword = qmap[0]
	}

	if len(startmap) == 0 {
		start = "0"
	} else {
		start = startmap[0]
	}

	if len(countmap) == 0 {
		count = "20"
	} else {
		count = countmap[0]
	}

	booklist, err := dbQueryBookListFromKeyword(keyword, start, count)
	checkErr(err)

	// if keyword search result is null, return 404
	if booklist.Total == 0 {
		var statusMsg ResourceStatus = *ResourceNotFound(req)
		w.WriteHeader(404)
		w.WriteJson(&statusMsg)
		return
	} else {
		w.WriteJson(&booklist)
		return
	}
}
