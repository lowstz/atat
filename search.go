package main

import (
	"github.com/garyburd/redigo/redis"	
	"github.com/huichen/sego"
	"time"
	"runtime"
//	"strings"
//	"fmt"
)


type QueryResult struct {
	total int
	books []string
}

//
func (engine *Engine) Query (keywords []string, start, count string) (QueryResult) {
	var queryResult QueryResult
	keywords = engine.keywordsFormat(keywords)
//	var cacheKey  = cacheResultPre + strings.Join(keywords, "+")
	cacheKey := cache.buildCacheKey(keywords)
	if cache.isIndexCacheKeyExist(cacheKey) {
		queryResult = cache.queryBooklistFromCache(cacheKey, start, count)
//		queryResult.total = cache.queryCountFromCache(cacheKey)
//		queryResult.books = cache.queryBooklistFromCache(cacheKey, start, count)
	} else {
		cache.storeCacheKey(keywords, cacheKey)
		queryResult = cache.queryBooklistFromCache(cacheKey, start, count)
		if queryResult.total > 0 {
			go cache.PreloadBookList(cacheKey)			
		}
//		queryResult.total = cache.queryCountFromCache(cacheKey)
//		queryResult.books = cache.queryBooklistFromCache(cacheKey, start, count)
//		fmt.Println(queryResult.books)
		// if len(queryResult.books) == 0 {
		// 	queryResult.books = append(queryResult.books, "00000000")
		// 	fmt.Println(queryResult.books)
		// }
	}
	return queryResult
}

func (engine *Engine) keywordsFormat (keywords []string) ([]string) {
	var newKeywords []string
	for _, keyword := range keywords {
		text := []byte(keyword)
		for _,word := range sego.SegmentsToSlice(engine.segmenter.Segment(text), false) {
			if !engine.st.IsStopToken(word) {
				newKeywords = append(newKeywords, word)
			}
		}
	}
	return newKeywords
}

func (engine *Engine) checkStatus() {
	
	for {
	rdb := cache.redispool.Get()
	key := indexConfigPre + "indexcomplete"
	status, err := redis.Int(rdb.Do("GET", key))
	checkErr(err)
		if status == 1 {
			controller.useCache = true
		} else {
			controller.useCache = false
		}
		runtime.GC()
		time.Sleep(30 * time.Second)
	}
}












