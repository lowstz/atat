package main

import (
	"log"
	"strconv"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)


// Return a *sql.DB object
func newDbConn() (*sql.DB, error) {
	database, err := sql.Open("mysql",
		config.Db.User+":"+config.Db.Password+"@tcp(["+
			config.Db.Host+"]:"+config.Db.Port+")/"+
			config.Db.Dbname+"?&charset=utf8")

	checkErr(err)
	database.SetMaxIdleConns(1)
	return database, err
}


// Statistics books number of search results.
func dbQueryBookCount(keyword string) (int, error) {
	db, err := newDbConn()
	checkErr(err)

	var itemCount int
	query := "select count(*) as count from bookinfo " + 
		"where book_title like '%" +
		keyword + "%'"

	err = db.QueryRow(query).Scan(&itemCount)
	checkErr(err)
	db.Close()
	return itemCount, err
}

// Fetch a single book information from book_id and 
// return as a Book struct.
func dbQueryBookFromBookId(book_id string) (Book, error) {
	db, err := newDbConn()
	checkErr(err)
	query := "select " +
		"id, book_id, book_title, book_author, " +
		"book_page_num, book_pubdate, book_publisher, " +
		"book_isbn, book_price, book_ref_no, book_category_num," +
		"book_summary from bookinfo where book_id = ?"

	rows, err := db.Query(query, book_id)
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
func dbQueryBookFromBookIsbn(book_isbn string) (Book, error) {
	db, err := newDbConn()
	checkErr(err)
	
	query := "select " +
		"id, book_id, book_title, book_author, " +
		"book_page_num, book_pubdate, book_publisher, " +
		"book_isbn, book_price, book_ref_no, book_category_num, book_summary " +
		" from bookinfo where book_isbn = ?"

	rows, err := db.Query(query, book_isbn)
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
func dbQueryBookListFromKeyword(keyword, start, count string) (BookList, error) {
	db, err := newDbConn()
	checkErr(err)

	total, err := dbQueryBookCount(keyword)
	checkErr(err)

	query := "select " +
		"book_id, book_isbn, book_title, book_author, book_ref_no" +
		" from bookinfo where book_title like '%" +
		keyword + "%' limit " +
		start + "," + count

	rows, err := db.Query(query)
	checkErr(err)

	var booklist BookList
	booklist.Start, _ = strconv.Atoi(start)
	booklist.Count, _ = strconv.Atoi(count)
	booklist.Total = total

	for rows.Next() {
		var book Book
		err := rows.Scan(
			&book.Book_id,
			&book.Book_isbn,
			&book.Book_title,
			&book.Book_author,
			&book.Book_ref_no,
		)
		if err != nil {
			log.Fatal(err)
		}
		booklist.Books = append(booklist.Books, book)
	}
	rows.Close()
	db.Close()
	return booklist, err
}
