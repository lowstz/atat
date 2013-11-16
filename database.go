package main

import (
//	"log"
	"strconv"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)



// Return a *sql.DB object
func newDbConn() (*sql.DB, error) {
	db, err := sql.Open("mysql",
		config.Db.User+":"+config.Db.Password+"@" + 
			config.Db.Protocol + "("+
			config.Db.Addr+")/"+
			config.Db.Dbname+"?&charset=utf8")

	checkErr(err)
	db.SetMaxIdleConns(1)
	return db, err
}


// Statistics books number of search results.
func dbQueryBookCount(keyword []string) (int, error) {
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
func dbQueryBookFromBookId(book_id string, fields []string) (Book, error) {
	db, err := newDbConn()
	checkErr(err)

	query := makeBookIdQuery(fields)
//	rows, err := db.Query(query)
	stmt, err := db.Prepare(query)
	checkErr(err)
	rows, err := stmt.Query(book_id)
	checkErr(err)
	var book Book
	for rows.Next() {
		err := rows.Scan(
			&book.Id,
			&book.Book_id,
			&book.Book_title,
			&book.Book_author,
			&book.Book_page_num,
			&book.Book_pubdate,
			&book.Book_publisher,
			&book.Book_isbn,
			&book.Book_price,
			&book.Book_ref_no,
			&book.Book_category_num,
			&book.Book_summary,
		)
		checkErr(err)
	}

	rows.Close()
	db.Close()
	return book, err
}

// Fetch a single book information from book_isbn and 
// return as a Book struct.
func dbQueryBookFromBookIsbn(book_isbn string, fields []string) (Book, error) {
	db, err := newDbConn()
	checkErr(err)
	
	query := makeBookIsbnQuery(fields)
	stmt, err := db.Prepare(query)
	checkErr(err)
	rows, err := stmt.Query(book_isbn)
	checkErr(err)

	var book Book
	for rows.Next() {
		err := rows.Scan(
			&book.Id,
			&book.Book_id,
			&book.Book_title,
			&book.Book_author,
			&book.Book_page_num,
			&book.Book_pubdate,
			&book.Book_publisher,
			&book.Book_isbn,
			&book.Book_price,
			&book.Book_ref_no,
			&book.Book_category_num,
			&book.Book_summary,
		)
		checkErr(err)
	}
	rows.Close()
	db.Close()
	return book, err
}

// Fetch a booklist from keyword and 
// return as a BookList struct.
func dbQueryBookListFromKeyword(keyword, fields []string, start, count string) (BookList, error) {
	db, err := newDbConn()
	checkErr(err)

	total, err := dbQueryBookCount(keyword)
	checkErr(err)

	query := makeKeywordSearchQuery(keyword, fields, start, count)
	rows, err := db.Query(query)
	checkErr(err)

	var booklist BookList
	booklist.Start, _ = strconv.Atoi(start)
	booklist.Count, _ = strconv.Atoi(count)
	booklist.Total = total

	for rows.Next() {
		var book Book
		// err := rows.Scan(
		// 	&book.Book_id,
		// 	&book.Book_isbn,
		// 	&book.Book_title,
		// 	&book.Book_author,
		// 	&book.Book_ref_no,
		// )
		err := rows.Scan(
			&book.Id,
			&book.Book_id,
			&book.Book_title,
			&book.Book_author,
			&book.Book_page_num,
			&book.Book_pubdate,
			&book.Book_publisher,
			&book.Book_isbn,
			&book.Book_price,
			&book.Book_ref_no,
			&book.Book_category_num,
			&book.Book_summary,
		)
		checkErr(err)
		booklist.Books = append(booklist.Books, book)
	}
	rows.Close()
	db.Close()
	return booklist, err
}
