package main

import (
	//	"fmt"
	"github.com/ant0ine/go-json-rest"
	"net/http"
	"strconv"
	"strings"
)

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

// controller for route /book/:id
// Query a book_id and response one single book information in a json.
func GetBookFromBookId(w *rest.ResponseWriter, req *rest.Request) {
	bookId := req.PathParam("id")
	parameter := req.URL.Query()
	fieldsMap := parameter["fields"]

	if len(bookId) != 8 {
		NotFoundError(w, req)
		return
	}

	fields := fieldsFilter(fieldsMap)
	//	checkErr(err)
	book, err := modelQueryBookFromBookId(bookId, fields)
	checkErr(err)

	if book.store_id != 0 {
		w.WriteJson(&book)
		return
	} else {
		NotFoundError(w, req)
		return
	}
}

// controller for route /book/isbn/:isbn
// Query from book_isbn and response one single book information in a json.
func GetBookFromBookISBN(w *rest.ResponseWriter, req *rest.Request) {
	book_isbn := req.PathParam("isbn")
	parameter := req.URL.Query()
	fieldsMap := parameter["fields"]

	if !isValidIsbn13(book_isbn) {
		NotFoundError(w, req)
		return
	}

	fields := fieldsFilter(fieldsMap)

	book, err := modelQueryBookFromBookIsbn(book_isbn, fields)
	checkErr(err)

	if book.store_id != 0 {
		w.WriteJson(&book)
		return
	} else {
		NotFoundError(w, req)
		return
	}
}

// controller for route /book/search/
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
		keyword = strings.Split(keywordMap[0], " ")
		keyword = replaceSpecialChar(keyword)
		if len(keyword) == 0 {
			NotFoundError(w, req)
			return
		}
	}

	fields = fieldsFilter(fieldsMap)

	if len(startMap) == 0 {
		start = "0"
	} else {
		start = startMap[0]
	}

	if len(countMap) == 0 {
		count = "20"
	} else {
		count = countMap[0]
		if countInt, _ := strconv.Atoi(count); countInt > 100 {
			count = "100"
		}
	}

	//  fmt.Println(keyword)
	//	fmt.Println(fields)
	//	fmt.Println(start)
	//	fmt.Println(count)
	booklist, err := modelQueryBookListFromKeyword(keyword, fields, start, count)
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

// return 404 json response
func NotFoundError(w *rest.ResponseWriter, req *rest.Request) {
	var statusMsg ResourceStatus = *ResourceNotFound(req)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.WriteJson(&statusMsg)
}

// Filter invalid field in the original fieldlist.
func fieldsFilter(fieldsMap []string) []string {
	var fields, secureFields []string
	if len(fieldsMap) == 0 {
		return defaultQueryFields
	} else {
		fields = strings.Split(fieldsMap[0], ",")
		for _, value := range fields {
			for _, secureValue := range defaultQueryFields {
				if value == secureValue {
					secureFields = append(secureFields, value)
				}
			}
		}
	}
	return secureFields
}
