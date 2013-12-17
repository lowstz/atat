package main

import (
	"github.com/ant0ine/go-json-rest"
	"log"
	"net/http"
	"os"
	"runtime"
//	"runtime/pprof"
//	_ "net/http/pprof"
//	"flag"
)

//var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {

	// flag.Parse()
    // if *cpuprofile != "" {
    //     f, err := os.Create(*cpuprofile)
    //     if err != nil {
    //         log.Fatal(err)
    //     }
	// 	pprof.StartCPUProfile(f)
    //     defer pprof.StopCPUProfile()
	// }
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()
	runtime.GOMAXPROCS(config.server.cpuCore)
	parseConfig("./conf/config.conf")
	logFile := getLogFile(config.global.logFile)
	startLogger := serverStartInfoLog(logFile)
	indexLogger := serverIndexInfoLog(logFile)
	startLogger.Println("AT-AT version", version)
	startLogger.Println(ascii_icon)
	startLogger.Println("Start atat server......")
	startLogger.Println("Server PID: ", os.Getpid())

	startLogger.Println("Starting Initialize Mysql Connection Pool")
	model.Init()
	startLogger.Println("Initialize Mysql Connection Pool Successful.")

	if config.global.cacheEable {
		startLogger.Println("Starting Initialize Redis Connection Pool.....")
		cache.Init()
		startLogger.Println("Initialize Redis Connection Pool Successful.")
		indexLogger.Println("Starting Initialize Indexer.....")
		engine.Init()
		go engine.IndexAll()
		go engine.checkStatus()
	}
//	startLogger.Println("Starting Initialize Controller")
//	controller.Init()
//	startLogger.Println("Initialize Controller Successful")

//	indexLogger.Println("Initialize Indexer Successful.")

	handler := rest.ResourceHandler{}
	handler.EnableGzip = config.server.gzip
	handler.EnableStatusService = config.server.gzip
	handler.EnableResponseStackTrace = config.global.debug
	handler.DisableJsonIndent = config.server.jsonIndent
	handler.Logger = log.New(logFile, "[Request] ", log.LstdFlags)

	handler.SetRoutes(
		rest.Route{"GET", "/book/search", controller.GetBookListFromKeyword},
		rest.Route{"HEAD", "/book/search", controller.GetBookListFromKeyword},
		rest.Route{"GET", "/book/:id", controller.GetBookFromBookId},
		rest.Route{"HEAD", "/book/:id", controller.GetBookFromBookId},
		rest.Route{"GET", "/book/isbn/:isbn", controller.GetBookFromBookISBN},
		rest.Route{"HEAD", "/book/isbn/:isbn", controller.GetBookFromBookISBN},
	)

	if config.server.listenAddr != "" && config.server.port != "" {
		ListenAddrPort := config.server.listenAddr + ":" + config.server.port
		startLogger.Println("Server listen on: ", ListenAddrPort)
		http.ListenAndServe(ListenAddrPort, &handler)

	} else {
		startLogger.Println("Server listen on: ", defaultListenAddrPort)
		http.ListenAndServe(defaultListenAddrPort, &handler)
	}
}
