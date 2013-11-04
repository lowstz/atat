package main

import (
//	"fmt"
	"strconv"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
    // "github.com/ziutek/mymysql/mysql"
    // _ "github.com/ziutek/mymysql/native" // Native engine
)




func newDbConn() (*sql.DB, error) {
	db, err := sql.Open("mysql", 
		config.Db.User + ":" + config.Db.Password + "@tcp(["+
		config.Db.Host + "]:" + config.Db.Port + ")/"+ 
		config.Db.Dbname + "?&charset=utf8")
	//	return db, err := sql.Open("mysql", config.Db.User+":"config.Db.Password+"@/"+config.Db.Dbname)

	if err != nil {
		panic( err )
	}
	return db, err
}

func dbGetBookCount() (int, error) {
	db, err := newDbConn()
	checkErr(err)
	defer db.Close()
	var itemCount int
	err = db.QueryRow("select count(*) as count from bookinfo").Scan(&itemCount)
	checkErr(err)
	defer db.Close()
	return itemCount, err
}

func dbGetBookFromBookId(book_id int) (Book, error){
	db, err := newDbConn()
	checkErr(err)
	defer db.Close()
	query := "select id, book_id, book_isbn, book_title" +
		" from bookinfo where book_id = ?"
	rows, err := db.Query(query, strconv.Itoa(book_id))
	checkErr(err)
	defer rows.Close()
	var book Book
	for rows.Next() {
		err := rows.Scan(
			&book.Id,
			&book.Book_id, 
			&book.Book_isbn, 
			&book.Book_title)
		checkErr(err)
	}
	return book, err
}
