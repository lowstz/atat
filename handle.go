package main

import (
	//	"log"
	"github.com/ant0ine/go-json-rest"
	"net/http"
	"strings"
)

// Store single book infomation.
type Book struct {
	Id                int    `json:"-"`
	Book_id           string `json:"id",omitempty`
	Book_title        string `json:"title,omitempty"`
	Book_author       string `json:"author,omitempty"`
	Book_category_num string `json:"category_num,omitempty"`
	Book_isbn         string `json:"isbn,omitempty"`
	Book_page_num     string `json:"page_num,omitempty"`
	Book_price        string `json:"price,omitempty"`
	Book_pubdate      string `json:"pubdate,omitempty"`
	Book_publisher    string `json:"publisher,omitempty"`
	Book_ref_no       string `json:"ref_no,omitempty"`
	Book_summary      string `json:"summary,omitempty"`
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
		Msg:     "Resource not found",
		Code:    404,
		Request: req.Method + " " + req.URL.Path + req.URL.RawQuery,
	}
}

// Handler for route /book/:id
// Query a book_id and response one single book information in a json.
func GetBookFromBookId(w *rest.ResponseWriter, req *rest.Request) {
	//	params := map[string][]string{"userId":[]string{ req.PathParam("id") }}
	//  book_id := req.PathParam("id")
	//	query := req.URL.Query()
	//	fields := query["fields"][0]
	//	fmt.Println(query)
	bookId := req.PathParam("id")
	parameter := req.URL.Query()
	fieldsMap := parameter["fields"]
	var fields []string
	if len(fieldsMap) == 0 {
		fields = []string{}
	} else {
		fields = strings.Split(fieldsMap[0], ",")
	}

	if len(bookId) != 8 {
		NotFoundError(w, req)
		return
	}
	//	checkErr(err)
	book, err := dbQueryBookFromBookId(bookId, fields)
	checkErr(err)

	if book.Id != 0 {
		w.WriteJson(&book)
		return
	} else {
		NotFoundError(w, req)
		return
	}
}

// Handler for route /book/isbn/:isbn
// Query from book_isbn and response one single book information in a json.
func GetBookFromBookISBN(w *rest.ResponseWriter, req *rest.Request) {
	//	params := map[string][]string{"userId":[]string{ req.PathParam("id") }}
	book_isbn := req.PathParam("isbn")
	parameter := req.URL.Query()
	fieldsMap := parameter["fields"]

	if !isValidIsbn13(book_isbn) {
		NotFoundError(w, req)
		return
	}

	var fields []string
	if len(fieldsMap) == 0 {
		fields = []string{}
	} else {
		fields = strings.Split(fieldsMap[0], ",")
	}

	book, err := dbQueryBookFromBookIsbn(book_isbn, fields)
	checkErr(err)

	if book.Id != 0 {
		w.WriteJson(&book)
		return
	} else {
		NotFoundError(w, req)
		return
	}
}

// Handler for route /book/search/
// Query from keyword and response a booklist with all book
// match keyword in a json.
func GetBookListFromKeyword(w *rest.ResponseWriter, req *rest.Request) {
	parameter := req.URL.Query()
	keywordMap := parameter["q"]
	startMap := parameter["start"]
	countMap := parameter["count"]
	fieldsMap := parameter["fields"]

	var (
		keyword, fields []string
		start, count    string
	)
	if len(keywordMap) == 0 {
		NotFoundError(w, req)
		return
	} else if keywordMap[0] == "" {
		NotFoundError(w, req)
		return
	} else {
		//		log.Println("keyword: ", keywordMap[0])
		keyword = strings.Split(keywordMap[0], " ")
	}

	if len(fieldsMap) == 0 {
		fields = []string{}
	} else {
		fields = strings.Split(fieldsMap[0], ",")
	}

	if len(startMap) == 0 {
		start = "0"
	} else {
		start = startMap[0]
	}

	if len(countMap) == 0 {
		count = "20"
	} else {
		count = countMap[0]
	}
	// log.Println(keyword)
	// log.Println(fields)
	// log.Println(start)
	// log.Println(count)
	booklist, err := dbQueryBookListFromKeyword(keyword, fields, start, count)
	checkErr(err)

	// if keyword search result is null, return 404
	if booklist.Total == 0 {
		NotFoundError(w, req)
		return
	} else {
		w.WriteJson(&booklist)
		return
	}
}

func NotFoundError(w *rest.ResponseWriter, req *rest.Request) {
	var statusMsg ResourceStatus = *ResourceNotFound(req)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.WriteJson(&statusMsg)
}
