package main

import (
	"fmt"
	"log"
	"runtime"
	"net/http"
	"github.com/ant0ine/go-json-rest"
)

func main() {
	parseConfig("./conf/config.conf")
	printVersion()
	fmt.Println("Start atat server......")
	printAsciiIcon()
	runtime.GOMAXPROCS(config.Server.CPUCore)
	handler := rest.ResourceHandler{}
	handler.EnableGzip = config.Server.Gzip
	handler.EnableStatusService = config.Server.Gzip
	handler.EnableResponseStackTrace = config.Debug
	handler.SetRoutes(
		rest.Route{"GET", "/book/:id", GetBookFromBookId},
		rest.Route{"GET", "/book/isbn/:isbn", GetBookFromBookISBN},
		rest.Route{"GET", "/book/search/", GetBookListFromKeyword},
	)
	if config.Server.ListenAddr != "" && config.Server.Port != "" {
		ListenAddrPort := config.Server.ListenAddr + ":" + config.Server.Port
		log.Println("Server listen on: ", ListenAddrPort)
		http.ListenAndServe(ListenAddrPort, &handler)

	}else {		
		log.Println("Server listen on: ", defaultListenAddrPort)
		http.ListenAndServe(defaultListenAddrPort, &handler)
	}
}





