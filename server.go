package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/ant0ine/go-json-rest"
)

func main() {
	parseConfig("./conf/config.conf")
	printVersion()
	fmt.Println("Start atat server......")
	printAsciiIcon()
	handler := rest.ResourceHandler{}
	handler.SetRoutes(
		rest.Route{"GET", "/book/:id", GetBookFromBookId},
		rest.Route{"GET", "/book/isbn/:isbn", GetBookFromBookISBN},
		rest.Route{"GET", "/book/search/", GetBookListFromKeyword},
	)
	if config.Server.ListenAddr != "" && config.Server.Port != "" {
		ListenAddrPort := config.Server.ListenAddr + ":" + config.Server.Port
		log.Println(ListenAddrPort)
		http.ListenAndServe(ListenAddrPort, &handler)

	}else {		
		log.Println(defaultListenAddrPort)
		http.ListenAndServe(defaultListenAddrPort, &handler)
	}
}
