package main

import (
	"github.com/ant0ine/go-json-rest"
	"log"
	"net/http"
	"os"
	"runtime"
)

func main() {
	parseConfig("./conf/config.conf")
	logFile := getLogFile(config.LogFile)
	startLogger := serverStartInfoLog(logFile)
	startLogger.Println("AT-AT version", version)
	startLogger.Println(ascii_icon)
	startLogger.Println("Start atat server......")
	startLogger.Println("Server PID: ", os.Getpid())

	runtime.GOMAXPROCS(config.Server.CPUCore)

	handler := rest.ResourceHandler{}
	handler.EnableGzip = config.Server.Gzip
	handler.EnableStatusService = config.Server.Gzip
	handler.EnableResponseStackTrace = config.Debug
	handler.DisableJsonIndent = config.Server.JsonIndent
	handler.Logger = log.New(logFile, "[Request] ", log.LstdFlags)

	handler.SetRoutes(
		rest.Route{"GET", "/book/search", GetBookListFromKeyword},
		rest.Route{"HEAD", "/book/search", GetBookListFromKeyword},
		rest.Route{"GET", "/book/:id", GetBookFromBookId},
		rest.Route{"HEAD", "/book/:id", GetBookFromBookId},
		rest.Route{"GET", "/book/isbn/:isbn", GetBookFromBookISBN},
		rest.Route{"HEAD", "/book/isbn/:isbn", GetBookFromBookISBN},
	)

	if config.Server.ListenAddr != "" && config.Server.Port != "" {
		ListenAddrPort := config.Server.ListenAddr + ":" + config.Server.Port
		startLogger.Println("Server listen on: ", ListenAddrPort)
		http.ListenAndServe(ListenAddrPort, &handler)

	} else {
		startLogger.Println("Server listen on: ", defaultListenAddrPort)
		http.ListenAndServe(defaultListenAddrPort, &handler)
	}
}
