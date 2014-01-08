package main

import (
	"github.com/ant0ine/go-json-rest"
	"net/http"
	"strconv"
	"strings"
//	"fmt"
)

// Use for return status if some exception occurs.
type Controller struct {
	useCache bool
}

type ResourceStatus struct {
	Msg     string
	Code    int
	Request string
}


var controller Controller

// func (controller *Controller) Init () {
// 	if config.global.cacheEable {
// 		if engine.indexComplete {
// 			controller.useCache = true
// 		}
// 	} else {
// 		controller.useCache = false
// 	}
// }

func (controller *Controller) GetAppRelease(w *rest.ResponseWriter, req *rest.Request) {
	parameter := req.URL.Query()
	verMap := parameter["ver"]
	var verCode string

	if len(verMap) == 0 {
		verCode = "latest"
	} else {
		if stringIsDigit(verMap[0]) {
			verCode = verMap[0]
		} else {
			verCode = "latest"
		}
	}
	appRelease := model.QueryAppRelease(verCode)
	if appRelease.VerName != ""  {
		w.WriteJson(&appRelease)
		return
	} else {
		controller.NotFoundError(w, req)
		return
	}
}


// controller for route /book/:id
// Query a book_id and response one single book information in a json.
func (controller *Controller) GetBookFromBookId(w *rest.ResponseWriter, req *rest.Request) {
	bookId := req.PathParam("id")
	parameter := req.URL.Query()
	fieldsMap := parameter["fields"]

	if len(bookId) != 8 {
		controller.NotFoundError(w, req)
		return
	}

	fields := fieldsFilter(fieldsMap)
	//	checkErr(err)
	book, err := model.QueryBookFromBookId(bookId, fields)
	checkErr(err)

	if book.Id != " " {
		// expire := time.Now().AddDate(0, 0, 1)
		// cookie := http.Cookie{Name: "testcookiename", Value: "testcookievalue", Path: "/", Expires: expire, MaxAge: 86400}
		// http.SetCookie(w, &cookie)
		w.WriteJson(&book)
		return
	} else {
		controller.NotFoundError(w, req)
		return
	}
}

// controller for route /book/isbn/:isbn
// Query from book_isbn and response one single book information in a json.
func (controller *Controller) GetBookFromBookISBN(w *rest.ResponseWriter, req *rest.Request) {
	book_isbn := req.PathParam("isbn")
	parameter := req.URL.Query()
	fieldsMap := parameter["fields"]

	if !isValidIsbn13(book_isbn) {
		controller.NotFoundError(w, req)
		return
	}

	fields := fieldsFilter(fieldsMap)

	book, err := model.QueryBookFromBookIsbn(book_isbn, fields)
	checkErr(err)

	if book.Isbn != " " {
		w.WriteJson(&book)
		return
	} else {
		controller.NotFoundError(w, req)
		return
	}
}

// controller for route /book/search/
// Query from keyword and response a booklist with all book
// match keyword in a json.
func (controller *Controller) GetBookListFromKeyword(w *rest.ResponseWriter, req *rest.Request) {
	parameter := req.URL.Query()
//	log.Println(req.URL.RawQuery)
	keywordMap := parameter["q"]
	startMap := parameter["start"]
	countMap := parameter["count"]
	fieldsMap := parameter["fields"]

	w.Header().Set("Access-Control-Allow-Origin", "*")
	var (
		keywords, fields []string
		start, count    string
	)
	if len(keywordMap) == 0 {
		controller.NotFoundError(w, req)
		return
	} else if keywordMap[0] == " " {
		controller.NotFoundError(w, req)
		return
	} else if keywordMap[0] == "" {
		controller.NotFoundError(w, req)
		return
	} else {
		keywords = strings.Split(keywordMap[0], " ")
		keywords = replaceSpecialChar(keywords)
		keywords = replaceUnclearChar(keywords)
		if len(keywords) == 0 {
			controller.NotFoundError(w, req)
			return
		}
	}

	fields = fieldsFilter(fieldsMap)
//	fmt.Println(fields)
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
	if controller.useCache {
		cacheResult := engine.Query(keywords, start, count)
		if cacheResult.total == 0 {
			controller.NotFoundError(w, req)
			return
		} else if len(cacheResult.books) == 0 {
			booklist := model.QueryBookListIfEmptySet(cacheResult, fields, start, count)
			w.WriteJson(&booklist)
			return			
		} else {
			booklist, err := model.QueryBookListFromCache(cacheResult, fields, start, count)
//			fmt.Println(booklist)
			checkErr(err)
			w.WriteJson(&booklist)
			return
		}
		// if keyword search result is null, return 404
	} else {
		booklist, err := model.QueryBookListFromKeyword(keywords, fields, start, count)		
		checkErr(err)
		
		// if keyword search result is null, return 404
		if booklist.Total == 0 {
			controller.NotFoundError(w, req)
			return
		} else {
			w.WriteJson(&booklist)
			return
		}
	}
}

// return 404 json response
func (controller *Controller)NotFoundError(w *rest.ResponseWriter, req *rest.Request) {
	var statusMsg ResourceStatus = *ResourceNotFound(req)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.WriteJson(&statusMsg)
}

// Return a ResourceStatus struct if Resource not found.
func ResourceNotFound(req *rest.Request) *ResourceStatus {
	return &ResourceStatus{
		Msg:     "Resource not found",
		Code:    404,
		Request: req.Method + " " + req.URL.Path + req.URL.RawQuery,
	}
}

// Filter invalid field in the original fieldlist.
func fieldsFilter(fieldsMap []string) []string {
	var fields, secureFields []string
	if len(fieldsMap) == 0 {
		return defaultQueryFields
	} else {
		fields = strings.Split(fieldsMap[0], ",")
		for _, value := range fields {
			if value == "douban_rate" {
				value = "average"
			}
			for _, secureValue := range defaultQueryFields {
				if value == secureValue {
					secureFields = append(secureFields, value)
				}
			}
		}
	}
	return secureFields
}
