package main

import (
	"fmt"
	"strings"
)

var (
	defaultQueryFields = []string{
		//		"store_id",
		"id",
		"title",
		"author",
		"page_num",
		"pubdate",
		"publisher",
		"isbn",
		"price",
		"ref_no",
		"category_num",
		"summary",
		"image",
		"images",
	}
)

// Make a sql query statement for /book/:id api.
func makeBookCountQuery(keyword []string) string {
	var condition []string
	selectQuery := "select count(*) as count from bookinfo where "

	for _, value := range keyword {
		bookCondition := fmt.Sprintf(
			"(author like '%%%s%%' or title like '%%%s%%' or isbn like '%%%s%%')",
			value, value, value)
		condition = append(condition, bookCondition)
	}
	conditionQuery := strings.Join(condition, " and ")
	sqlQuery := selectQuery + conditionQuery
	return sqlQuery
}

// Make a sql query statement for /book/:id api.
func makeBookIdQuery(fields []string) string {
	if len(fields) == 0 {
		fields = defaultQueryFields
	}
	selectQuery := "select "
	fieldsQuery := strings.Join(fields, ",")
	conditionQuery := " from bookinfo where id=?"
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
	conditionQuery := " from bookinfo where isbn=?"
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
			"(author like '%%%s%%' or title like '%%%s%%' or isbn like '%%%s%%')",
			value, value, value)
		condition = append(condition, bookCondition)
	}
	conditionQuery := strings.Join(condition, " and ")
	pageQuery := " limit " + start + "," + count

	sqlQuery := selectQuery + fieldsQuery +
		fromQuery + conditionQuery + pageQuery
	return sqlQuery
}

// Replace unsafe request keyword query parameters.
func replaceSpecialChar(keyword []string) []string {
	var secureKeyword []string
	for _, value := range keyword {
		switch value {
		case "<":
			break
		case ">":
			break
		case "*":
			break
		case "_":
			break
		case "?":
			break
		case "'":
			break
		case ",":
			break
		case "/":
			break
		case ":":
			break
		case "*/":
			break
		case "\r\n":
			break
		default:
			secureKeyword = append(secureKeyword, value)
		}
	}
//	fmt.Println(secureKeyword)
	return secureKeyword
}
