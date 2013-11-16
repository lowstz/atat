package main

import (
	"fmt"
	"strings"
)

var (
	defaultQueryFields = []string{
		"id",
		"book_id",
		"book_title",
		"book_author",
		"book_page_num",
		"book_pubdate",
		"book_publisher",
		"book_isbn",
		"book_price",
		"book_ref_no",
		"book_category_num",
		"book_summary ",
	}
)

// Make a sql query statement for /book/:id api.
func makeBookCountQuery(keyword []string) string {
	var condition []string
	selectQuery := "select count(*) as count from bookinfo where "
	for _, value := range keyword {
		bookCondition := fmt.Sprintf(
			"(book_author like '%%%s%%' or book_title like '%%%s%%')", value, value)
		condition = append(condition, bookCondition)
	}
	conditionQuery := strings.Join(condition, " and ")
	sqlQuery :=  selectQuery + conditionQuery
	return sqlQuery
}

// Make a sql query statement for /book/:id api.
func makeBookIdQuery(fields []string) string {
	if len(fields) == 0 {
		fields = defaultQueryFields
	}
	selectQuery := "select "
	fieldsQuery := strings.Join(fields, ",")
	conditionQuery := " from bookinfo where book_id=?"
//	conditionQuery := fmt.Sprintf("book_id='%s'", book_id)

	sqlQuery := selectQuery + fieldsQuery + conditionQuery
	return sqlQuery
}

// Make a sql query statement for /book/isbn/:isbn api.
func makeBookIsbnQuery(fields []string) string {
	if len(fields) == 0 {
		fields = defaultQueryFields
	}
	selectQuery := "select "
	fieldsQuery := strings.Join(fields, ",")
	conditionQuery := " from bookinfo where book_isbn=?"
//	conditionQuery := fmt.Sprintf("book_isbn='%s'", book_isbn)

	sqlQuery := selectQuery + fieldsQuery + conditionQuery
	return sqlQuery
}

// Make a sql query statement for /book/search api.
func makeKeywordSearchQuery(keyword, fields []string, start, count string) string {
	if len(fields) == 0 {
		fields = defaultQueryFields
	}
	var condition []string
	selectQuery := "select "
	fieldsQuery := strings.Join(fields, ",")
	fromQuery := " from bookinfo where "
	for _, value := range keyword {
		bookCondition := fmt.Sprintf(
			"(book_author like '%%%s%%' or book_title like '%%%s%%' or book_isbn='%s')",
			value, value, value)
		condition = append(condition, bookCondition)
	}
	conditionQuery := strings.Join(condition, " and ")
	pageQuery := " limit " + start + "," + count

	sqlQuery := selectQuery + fieldsQuery +
		fromQuery + conditionQuery + pageQuery
	return sqlQuery
}
