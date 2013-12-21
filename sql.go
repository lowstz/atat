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
func makeBookIdQuery(fields []string) string {
	if len(fields) == 0 {
		fields = defaultQueryFields
	}
	selectQuery := "select "
	fieldsQuery := strings.Join(fields, ",")
	//	conditionQuery := " from  book_items left join book_category on book_items.category_num=book_category.category_num where id=?"
	conditionQuery := " from book_items " +
		"left join book_category on book_items.category_num=book_category.category_num where id=?"
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
	conditionQuery := " from book_items left join book_category on book_items.category_num=book_category.category_num where isbn=?"
	// conditionQuery := " from book_items left join book_douban on book_items.id=book_douban.book_id " +
	//	"left jowhere isbn=?"

	//	conditionQuery := fmt.Sprintf("book_isbn='%s'", book_isbn)

	sqlQuery := selectQuery + fieldsQuery + conditionQuery
//	fmt.Println(sqlQuery)
	return sqlQuery
}

// Make a sql query statement for book number count api.
func makeBookCountQuery(keyword []string) string {
	var condition []string
	selectQuery := "select count(*) as count from book_items where "

	for _, value := range keyword {
		bookCondition := fmt.Sprintf(
			"(author like '%%%s%%' or title like '%%%s%%' or summary like '%%%s%%' or isbn='%s')",
			value, value, value, value)
		condition = append(condition, bookCondition)
	}
	conditionQuery := strings.Join(condition, " and ")
	sqlQuery := selectQuery + conditionQuery
//	fmt.Println(sqlQuery)
//	sqlHash := Sha1Hasher(sqlQuery)
	return sqlQuery  // sqlHash
}

// Make a sql query statement for /book/search api.
func makeKeywordSearchQuery(keywords, fields []string, start, count string) string {
	if len(fields) == 0 {
		fields = defaultQueryFields
	}

	var condition []string
	selectQuery := "select "
	fieldsQuery := strings.Join(fields, ",")
	//	oldfromQuery := " from book_items left join book_category on book_items.category_num=book_category.category_num where "
	fromQuery := " from book_items " +
		"left join book_category on book_items.category_num=book_category.category_num where "
	for _, value := range keywords {
		bookCondition := fmt.Sprintf(
			"(author like '%%%s%%' or title like '%%%s%%' or summary like '%%%s%%' or isbn='%s' )",
			value, value, value, value)
		condition = append(condition, bookCondition)
	}

	conditionQuery := strings.Join(condition, " and ")
//	sortQuery := " order by (average*100+num_raters*0.3) desc "
	sortQuery := " order by rank desc "
	pageQuery := " limit " + start + "," + count

	sqlQuery := selectQuery + fieldsQuery +
		fromQuery + conditionQuery + sortQuery + pageQuery
//	fmt.Println(sqlQuery)
	return sqlQuery
}

func makeCacheKeywordSearchQuery(idlist , fields []string) string {
	if len(fields) == 0 {
		fields = defaultQueryFields
	}
	selectQuery := "select "
	fieldsQuery := strings.Join(fields, ",")
	condition   := strings.Join(idlist, ",")
	fromQuery   := " from book_items left join book_category on book_items.category_num=book_category.category_num where id in (" + condition + ") order by rank desc"
//	fromQuery   := " from book_items where id in (" + condition + ")"
	sqlQuery := selectQuery + fieldsQuery + fromQuery
//	fmt.Println(sqlQuery)
	return sqlQuery
}

func makeAppRelease(verCode string) string {
	if verCode == "latest" {
		sqlQuery := "select * from app_release order by vercode desc limit 0,1"
		return sqlQuery
	} else {
		sqlQuery := "select * from app_release where vercode=" + verCode
		return sqlQuery
	}
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
				clearKeyword = append(clearKeyword, "C语言")
				break
			case "c":
				clearKeyword = append(clearKeyword, "c语言")
				break
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


// func handleNullField(fields []string, field string) []string {
// 	for p, v := range fields {
// 		if v == field {
// 			fields[p] = "coalesce(" + field + ", 'unknown') as " + field
// 			//			return p
// 		}
// 	}
// 	return fields
// }
