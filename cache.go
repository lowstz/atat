package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
	"strings"
//	"fmt"
)

const (
	cacheResultPre  = "atat:cache:result:"
	cacheInstancePre = "atat:cache:instance:"
	cacheResultSort  = "atat:index:_score_:*"
	cacheResultTime  = 3600
)

//	redispool *redis.Pool
type RDB struct {
	redispool *redis.Pool
}

var cache RDB

//var redispool *redis.Pool
func (cache *RDB)Init() {
	cache.redispool = &redis.Pool{
		MaxIdle:     3,
		MaxActive:   0,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("unix", "/tmp/redis.sock")
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}


func (cache *RDB) isIndexCacheKeyExist (cacheKey string) (bool) {
	rdb := cache.redispool.Get()
	exist, err := redis.Int(rdb.Do("EXISTS", cacheKey))
	checkErr(err)
	rdb.Close()
	return exist == 1
}

func (cache *RDB) buildCacheKey(keywords []string) (string){
	return cacheResultPre + strings.Join(keywords, "|")
}

func (cache *RDB) buildStoreKey(keywords []string) (string) {
	return cacheInstancePre + strings.Join(keywords, "|")
}

func (cache *RDB) storeCacheKey(keywords []string, cacheKey string) {
	rdb := cache.redispool.Get()
	args := []interface{}{cacheKey}
	for _,keyword := range keywords {
		indexKey := indexKeywordPre + keyword
		args = append(args, indexKey)
	}
//	fmt.Println(args)
	rdb.Send("SINTERSTORE", args...)
	rdb.Send("EXPIRE", cacheKey, cacheResultTime)
	rdb.Flush()
	rdb.Close()
}


// func (cache *RDB) storeInstanceKey(keywords []string, cacheKey string) {
// //	rdb := cache.redispool.Get()
	
// }

// func (cache *RDB) queryCountFromCache(cacheKey string) (total int) {
// 	rdb := cache.redispool.Get()
// 	rdb.Send("SCARD", cacheKey)
// //	rdb.db.Send("SORT", cacheKey, "BY", cacheResultSort, "DESC", "LIMIT", start, count)
// 	rdb.Flush()
// 	total, err := redis.Int(rdb.Receive())
// 	checkErr(err)
// 	rdb.Close()
// 	return total
// }

func (cache *RDB) queryBooklistFromCache(cacheKey, start, count string) (QueryResult) {
	var queryResult QueryResult
	rdb := cache.redispool.Get()
	rdb.Send("SCARD", cacheKey)
//	rdb.db.Send("SORT", cacheKey, "BY", cacheResultSort, "DESC", "LIMIT", start, count)

	rdb.Send("SORT", cacheKey, "BY", cacheResultSort, "DESC", "LIMIT", start, count)
	rdb.Flush()
	total, err := redis.Int(rdb.Receive())
	checkErr(err)
	reply, err := redis.Values(rdb.Receive())
	checkErr(err)
	rdb.Close()
	var result []string
	for _, x := range reply {
		var v, ok = x.([]byte)
		if ok {
			result = append(result, string(v))
		}
	}
	queryResult.total = total
	queryResult.books = result
	return queryResult
}

func (cache *RDB) PreloadBookList(cacheKey string) {
	rdb := cache.redispool.Get()
	rdb.Send("SORT", cacheKey, "BY", cacheResultSort, "DESC")
	rdb.Flush()
	reply, err := redis.Values(rdb.Receive())
	checkErr(err)
	rdb.Close()
	var result []string
	var fields []string
	for _, x := range reply {
		var v, ok = x.([]byte)
		if ok {
			result = append(result, string(v))
		}
	}
	go model.PreloadBookList(result, fields)
}













