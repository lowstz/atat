package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strconv"
	"strings"
)

// Store single book infomation.
type Book struct {
	store_id     uint64 `json:"-"`
	Id           string `json:"id",omitempty`
	Title        string `json:"title,omitempty"`
	Author       string `json:"author,omitempty"`
	Category_num string `json:"category_num,omitempty"`
	Isbn         string `json:"isbn,omitempty"`
	Page_num     string `json:"page_num,omitempty"`
	Price        string `json:"price,omitempty"`
	Pubdate      string `json:"pubdate,omitempty"`
	Publisher    string `json:"publisher,omitempty"`
	Ref_no       string `json:"ref_no,omitempty"`
	Summary      string `json:"summary,omitempty"`
}

// Store multiple books.
type BookList struct {
	Start int
	Count int
	Total int
	Books []Book
}

// Return a *sql.DB object
func newDbConn() (*sql.DB, error) {
	db, err := sql.Open("mysql",
		config.Db.User+":"+config.Db.Password+"@"+
			config.Db.Protocol+"("+
			config.Db.Addr+")/"+
			config.Db.Dbname+"?&charset=utf8")

	checkErr(err)
	db.SetMaxIdleConns(100)
	return db, err
}

// Statistics books number of search results.
func modelQueryBookCount(keyword []string) (int, error) {
	db, err := newDbConn()
	checkErr(err)

	var itemCount int
	query := makeBookCountQuery(keyword)
	//	log.Println(query)
	err = db.QueryRow(query).Scan(&itemCount)
	checkErr(err)
	db.Close()
	return itemCount, err
}

// Fetch a single book information from book_id and
// return as a Book struct.
func modelQueryBookFromBookId(book_id string, fields []string) (Book, error) {
	db, err := newDbConn()
	checkErr(err)

	query := makeBookIdQuery(fields)
	//	rows, err := db.Query(query)
	stmt, err := db.Prepare(query)
	checkErr(err)
	rows, err := stmt.Query(book_id)
	checkErr(err)

	cols, _ := rows.Columns()
	buff := make([]interface{}, len(cols)) // tmp slice
	data := make([]string, len(cols))      // store slice
	for i, _ := range buff {
		buff[i] = &data[i]
	}

	for rows.Next() {
		err := rows.Scan(buff...)
		checkErr(err)
	}
	book := mapDataToStruct(fields, data)
	rows.Close()
	db.Close()
	return book, err
}

// Fetch a single book information from book_isbn and
// return as a Book struct.
func modelQueryBookFromBookIsbn(book_isbn string, fields []string) (Book, error) {
	db, err := newDbConn()
	checkErr(err)

	query := makeBookIsbnQuery(fields)
	stmt, err := db.Prepare(query)
	checkErr(err)
	rows, err := stmt.Query(book_isbn)
	checkErr(err)

	cols, _ := rows.Columns()
	buff := make([]interface{}, len(cols)) // tmp slice
	data := make([]string, len(cols))      // store slice
	for i, _ := range buff {
		buff[i] = &data[i]
	}

	for rows.Next() {
		err := rows.Scan(buff...)
		checkErr(err)
	}
	book := mapDataToStruct(fields, data)
	rows.Close()
	db.Close()
	return book, err
}

// Fetch a booklist from keyword and
// return as a BookList struct.
func modelQueryBookListFromKeyword(keyword, fields []string, start, count string) (BookList, error) {
	//	fmt.Println(fields)
	db, err := newDbConn()
	checkErr(err)

	total, err := modelQueryBookCount(keyword)
	checkErr(err)

	query := makeKeywordSearchQuery(keyword, fields, start, count)
	rows, err := db.Query(query)
	checkErr(err)

	var booklist BookList
	booklist.Start, _ = strconv.Atoi(start)
	booklist.Count, _ = strconv.Atoi(count)
	booklist.Total = total

	cols, _ := rows.Columns()
	buff := make([]interface{}, len(cols)) // tmp slice
	data := make([]string, len(cols))      // store slice
	for i, _ := range buff {
		buff[i] = &data[i]
	}

	for rows.Next() {
		err := rows.Scan(buff...)
		checkErr(err)
		//	fmt.Println(data)
		book := mapDataToStruct(fields, data)
		booklist.Books = append(booklist.Books, book)
	}
	rows.Close()
	db.Close()
	return booklist, err
}

// map Fetched data slice to Book struct
func mapDataToStruct(fields, data []string) Book {
	if len(fields) == 0 {
		fields = defaultQueryFields
	}
	var book Book
	bookValue := reflect.ValueOf(&book).Elem()
	//	bookType := reflect.TypeOf(&book).Elem()
	for key, value := range fields {
		kind := bookValue.FieldByName(strings.Title(value)).Kind()
		// fmt.Println(kind)
		// fmt.Println(book)
		switch kind {
		case reflect.Uint64:
			id, _ := strconv.Atoi(data[key])
			bookValue.FieldByName(strings.Title(value)).SetUint(uint64(id))
		case reflect.String:
			bookValue.FieldByName(strings.Title(value)).SetString(data[key])
		default:
		}
	}
	return book
}
