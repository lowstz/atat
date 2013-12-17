package main

import (
	"bufio"
//	"database/sql"
	"github.com/garyburd/redigo/redis"
//	_ "github.com/go-sql-driver/mysql"
	"github.com/huichen/sego"
	"log"
	"os"
	"strconv"
	"time"
	"github.com/huichen/pinyin"
)

const (
	indexKeywordPre    = "atat:index:keyword:"
	indexScorePre      = "atat:index:_score_:"
	indexConfigPre     = "atat:index:config:"
	indexInstancePre   = "atat:index:instance:"
	dictionaryFilePath = "./data/dictionary.txt"
	stopTokenFilePath  = "./data/stop_tokens.txt"
	pinyinTableFilePath  = "./data/pinyin_table.txt"
)

// var (
// 	segmenter = sego.Segmenter{}
// 	numThreads = runtime.NumCPU()
// 	task = make(chan BookIndex, numThreads)
// 	numRuns  = 20
// 	st StopTokens
// )
var engine Engine

type StopTokens struct {
	stopTokens map[string]bool
}

type Engine struct {
	segmenter     sego.Segmenter
	st            StopTokens
	py            pinyin.Pinyin
	numThreads    int
	task          chan BookEngine
	numRuns       int
	pinyinMatch   bool
	indexComplete bool
}

type BookEngine struct {
	Id      int
	Isbn    string
	Author  string
	Title   string
	Summary string
	Rank    float64
}

func (engine *Engine) Init() {
	// Load Dictionary
	engine.segmenter.LoadDictionary(ExpandPath(dictionaryFilePath))
	// Load StopToken
	engine.st.Init(ExpandPath(stopTokenFilePath))
	engine.py.Init(ExpandPath(pinyinTableFilePath))
	engine.numThreads = config.server.cpuCore
	engine.task = make(chan BookEngine, engine.numThreads)
	engine.numRuns = 20
	engine.pinyinMatch = true
	key := indexConfigPre + "indexcomplete"
	rdb := cache.redispool.Get()
	existindexStatus,err := redis.Int(rdb.Do("EXISTS", key))
	if existindexStatus == 0 {
		rdb.Do("SET", key, 0)
	}
	indexStatus, err := redis.Int(rdb.Do("GET", key))
	checkErr(err)
	if indexStatus == 1 {
		engine.indexComplete = true
	} else {
		engine.indexComplete = false
	}
	rdb.Close()
}

func (engine *Engine) IndexAll() {
	logFile := getLogFile(config.global.logFile)
	indexLogger := serverIndexInfoLog(logFile)

	if engine.indexComplete {
		indexLogger.Println("Don't repeat create index")
	} else {
		// Load database document to Memory
		query := "select id, isbn, author, title, summary, (average*100+num_raters*0.1) as newrank from book_items"
		rows, err := model.db.Query(query)
		checkErr(err)
		var booklist []BookEngine
		var book BookEngine
		size := 0

		for rows.Next() {
			err := rows.Scan(
				&book.Id,
				&book.Isbn,
				&book.Author,
				&book.Title,
				&book.Summary,
				&book.Rank)
			checkErr(err)
			size += len([]byte(book.Author))
			size += len([]byte(book.Title))
			size += len([]byte(book.Summary))
			booklist = append(booklist, book)
		}

		// Lanuch splitWorker goroutine
		for i := 0; i < engine.numThreads; i++ {
			go engine.splitWorker()
		}
		indexLogger.Println("分词开始")

		// record time
		t0 := time.Now()

		// Parallel to split the documents
		for i := 0; i < engine.numRuns; i++ {
			for _, value := range booklist {
				engine.task <- value
			}
		}

		// record time and calculate the speed of split documents
		t1 := time.Now()
		indexLogger.Printf("分词花费时间 %v", t1.Sub(t0))
		indexLogger.Printf("分词速度 %f MB/s", float64(size*engine.numRuns)/t1.Sub(t0).Seconds()/(1024*1024))

		// Set the index status
//		rdb := redispool.Get()
		rdb := cache.redispool.Get()
		key := indexConfigPre + "indexcomplete"
		rdb.Do("SET", key, 1)
		rdb.Close()
	}
}

func (engine *Engine) splitWorker() {
	for {
		rdb := cache.redispool.Get()
		book := <-engine.task
		socreKey := indexScorePre + strconv.Itoa(book.Id)
		rdb.Send("SET", socreKey, book.Rank)
		if  isValidIsbn13(book.Isbn) {
			isbnKey := indexKeywordPre + book.Isbn
			rdb.Send("SADD", isbnKey, book.Id)			
		}
		text := []byte(book.Author + "," + book.Title + "," + book.Summary)
		for _, word := range sego.SegmentsToSlice(engine.segmenter.Segment(text), true) {
			if !engine.st.IsStopToken(word) {
				key := indexKeywordPre + word
				rdb.Send("SADD", key, book.Id)
				if engine.pinyinMatch {
					pinyinWord := engine.pinyinConverter(word)
					if pinyinWord != "" {
						pinyinKey := indexKeywordPre + pinyinWord
						rdb.Send("SADD", pinyinKey, book.Id)						
					}
				}
			}
		}
		rdb.Flush()
		rdb.Close()
	}
}

func (engine *Engine) pinyinConverter(zhstr string) (string) {
	var pystr string
	for _,word := range zhstr {
		pystr += engine.py.GetPinyin(word, false)
	}
	return pystr
}

func (st *StopTokens) Init(stopTokenFile string) {
	st.stopTokens = make(map[string]bool)
	if stopTokenFile == "" {
		return
	}

	file, err := os.Open(ExpandPath(stopTokenFile))
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if text != "" {
			st.stopTokens[text] = true
		}
	}
}

func (st *StopTokens) IsStopToken(token string) bool {
	_, found := st.stopTokens[token]
	return found
}
