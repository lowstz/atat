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
		"category_name",
		"average",
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
			"(author like '%%%s%%' or title like '%%%s%%' or summary like '%%%s%%' or isbn='%s' )",
			value, value, value, value)
		condition = append(condition, bookCondition)
	}
	conditionQuery := strings.Join(condition, " and ")
	sqlQuery := selectQuery + conditionQuery
//	fmt.Println(sqlQuery)
	return sqlQuery
}

// Make a sql query statement for /book/:id api.
func makeBookIdQuery(fields []string) string {
	if len(fields) == 0 {
		fields = defaultQueryFields
	}
	selectQuery := "select "
	fieldsQuery := strings.Join(fields, ",")
	//	conditionQuery := " from  bookinfo left join book_category on bookinfo.category_num=book_category.category_num where id=?"
	conditionQuery := " from bookinfo left join douban on bookinfo.id=douban.book_id " +
		"left join book_category on bookinfo.category_num=book_category.category_num where id=?"
	//	conditionQuery := fmt.Sprintf("book_id='%s'", book_id)

	sqlQuery := selectQuery + fieldsQuery + conditionQuery
//	fmt.Println(sqlQuery)
	return sqlQuery
}

// Make a sql query statement for /book/isbn/:isbn api.
func makeBookIsbnQuery(fields []string) string {
	if len(fields) == 0 {
		fields = defaultQueryFields
	}
	selectQuery := "select "
	fieldsQuery := strings.Join(fields, ",")
	//	conditionQuery := " from bookinfo left join book_category on bookinfo.category_num=book_category.category_num where isbn=?"
	conditionQuery := " from bookinfo left join douban on bookinfo.id=douban.book_id " +
		"left join book_category on bookinfo.category_num=book_category.category_num where isbn=?"

	//	conditionQuery := fmt.Sprintf("book_isbn='%s'", book_isbn)

	sqlQuery := selectQuery + fieldsQuery + conditionQuery
//	fmt.Println(sqlQuery)
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
	//	oldfromQuery := " from bookinfo left join book_category on bookinfo.category_num=book_category.category_num where "
	fromQuery := " from bookinfo left join douban on bookinfo.id=douban.book_id " +
		"left join book_category on bookinfo.category_num=book_category.category_num where "
	for _, value := range keyword {
		bookCondition := fmt.Sprintf(
			"(author like '%%%s%%' or title like '%%%s%%' or summary like '%%%s%%' or isbn='%s' )",
			value, value, value, value)
		condition = append(condition, bookCondition)
	}

	conditionQuery := strings.Join(condition, " and ")
	sortQuery := " order by (average*100+num_raters*0.3) desc "
	pageQuery := " limit " + start + "," + count

	sqlQuery := selectQuery + fieldsQuery + 
		fromQuery + conditionQuery + sortQuery + pageQuery
//	fmt.Println(sqlQuery)
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


// replace unclear search keyword
// TODO: collect a unclear search keyword database.
func replaceUnclearChar(keyword []string) []string {
	var clearKeyword []string
	if len(keyword) == 1 {
		for _, value := range keyword {
			switch value {
			case "C":
				value = "C语言"
				clearKeyword = append(clearKeyword, value)
				break;
			case "c":
				value = "c语言"
				clearKeyword = append(clearKeyword, value)
				break;
			default:
				clearKeyword = append(clearKeyword, value)
			}
		}
	} else {
		clearKeyword = keyword
	}
	//	fmt.Println(secureKeyword)
	return clearKeyword
}
