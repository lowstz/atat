package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strconv"
	"strings"

//	"fmt"
)

// Store single book infomation.
type Book struct {
	store_id      uint64     `json:"-"`
	Id            string     `json:"id,omitempty"`
	Title         string     `json:"title,omitempty"`
	Author        string     `json:"author,omitempty"`
	Category_name string     `json:"category_name,omitempty"`
	Isbn          string     `json:"isbn,omitempty"`
	Page_num      string     `json:"page_num,omitempty"`
	Price         string     `json:"price,omitempty"`
	Pubdate       string     `json:"pubdate,omitempty"`
	Publisher     string     `json:"publisher,omitempty"`
	Ref_no        string     `json:"ref_no,omitempty"`
	Average       float64    `json:"douban_rate"`
	Summary       string     `json:"summary,omitempty"`
	Image         string     `json:"image,omitempty"`
	Images        Book_Image `json:"images,omitempty"`
}

type Book_Image struct {
	Large  string `json:"large,omitempty"`
	Medium string `json:"medium,omitempty"`
	Small  string `json:"small,omitempty"`
}

// Store multiple books.
type BookList struct {
	Start int
	Count int
	Total int
	Books []Book
}

type Model struct {
	db *sql.DB
}

var model Model

// Return a *sql.DB object
func (model *Model) Init() {
	db, err := sql.Open("mysql",
		config.db.user+":"+config.db.password+"@"+
			config.db.protocol+"("+
			config.db.addr+")/"+
			config.db.dbname+"?&charset=utf8")

	checkErr(err)
	db.SetMaxIdleConns(3)
	model.db = db
}

// Statistics books number of search results.
func (model *Model) QueryBookCount(keyword []string) (int, error) {

	var itemCount int
	query := makeBookCountQuery(keyword)
	//	log.Println(query)
	err := model.db.QueryRow(query).Scan(&itemCount)
	checkErr(err)
//	model.db.Close()
	return itemCount, err
}

// Fetch a single book information from book_id and
// return as a Book struct.
func (model *Model) QueryBookFromBookId(book_id string, fields []string) (Book, error) {

	query := makeBookIdQuery(fields)
	//	rows, err := db.Query(query)
	stmt, err := model.db.Prepare(query)
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
//	db.Close()
	return book, err
}

// Fetch a single book information from book_isbn and
// return as a Book struct.
func (model *Model) QueryBookFromBookIsbn(book_isbn string, fields []string) (Book, error) {

	query := makeBookIsbnQuery(fields)
	stmt, err := model.db.Prepare(query)
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
//	model.db.Close()
	return book, err
}

// Fetch a booklist from keyword and
// return as a BookList struct.
func (model *Model) QueryBookListFromKeyword(keywords, fields []string, start, count string) (BookList, error) {
	//	fmt.Println(fields)
	total, err := model.QueryBookCount(keywords)
	checkErr(err)

	query := makeKeywordSearchQuery(keywords, fields, start, count)
	rows, err := model.db.Query(query)
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
//	db.Close()
	return booklist, err
}

func (model *Model) QueryBookListFromCache(cacheResult QueryResult, fields []string, start, count string) (BookList, error) {
	var booklist BookList
	booklist.Start, _ = strconv.Atoi(start)
	booklist.Count, _ = strconv.Atoi(count)
	booklist.Total = cacheResult.total
	
	query := makeCacheKeywordSearchQuery(cacheResult.books, fields)
	rows, err := model.db.Query(query)
	checkErr(err)

	cols, _ := rows.Columns()
	buff := make([]interface{}, len(cols))
	data := make([]string, len(cols))
	for i, _ := range buff {
		buff[i] = &data[i]
	}

	for rows.Next() {
		err := rows.Scan(buff...)
		checkErr(err)
		book := mapDataToStruct(fields, data)
		booklist.Books = append(booklist.Books, book)
	}
	rows.Close()
	return booklist, err
}

func (model *Model) PreloadBookList(idlist, fields []string) {
	query := makeCacheKeywordSearchQuery(idlist, fields)
	rows, err := model.db.Query(query)
	checkErr(err)
	rows.Close()
}

func (model *Model) QueryBookListIfEmptySet (cacheResult QueryResult, fields []string, start, count string) (BookList) {
	var booklist BookList
	booklist.Start, _ = strconv.Atoi(start)
	booklist.Count, _ = strconv.Atoi(count)
	booklist.Total = cacheResult.total
	booklist.Books = []Book{}
	// query := makeCacheKeywordSearchQuery(cacheResult.books, fields)
	// rows, err := model.db.Query(query)
	// checkErr(err)

	// cols, _ := rows.Columns()
	// buff := make([]interface{}, len(cols))
	// data := make([]string, len(cols))
	// for i, _ := range buff {
	// 	buff[i] = &data[i]
	// }

	// for rows.Next() {
	// 	err := rows.Scan(buff...)
	// 	checkErr(err)
	// 	book := mapDataToStruct(fields, data)
	// 	booklist.Books = append(booklist.Books, book)
	// }
	// rows.Close()
	return booklist
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
		//		fmt.Println(kind)
		switch kind {
		case reflect.Uint64:
			id, _ := strconv.Atoi(data[key])
			bookValue.FieldByName(strings.Title(value)).SetUint(uint64(id))
		case reflect.String:
			//			fmt.Println(key, value)
			if len(data[key]) == 0 {
				data[key] = " "
			}
			bookValue.FieldByName(strings.Title(value)).SetString(data[key])
		case reflect.Float64:
			//			fmt.Println(key, value)
			if len(data[key]) == 0 {
				data[key] = " "
			}
//			fmt.Println(data[key])
			floatdata, err := strconv.ParseFloat(data[key],64)
			checkErr(err)
			bookValue.FieldByName(strings.Title(value)).SetFloat(floatdata)
		case reflect.Struct:
			f := make(map[string]string)
			json.Unmarshal([]byte(data[key]), &f)
			//			fmt.Println(value)
			//			fmt.Println(f)
			for img_type, url := range f {
				//				fmt.Println(img_type, url)
				bookValue.FieldByName(strings.Title(value)).FieldByName(strings.Title(img_type)).SetString(url)
			}
		default:
		}
	}
	return book
}
