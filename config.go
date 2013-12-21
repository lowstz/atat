package main

import (
	//	"log"
	"github.com/Unknwon/goconfig"
	"strconv"
)

const (
	version               = "0.3"
	defaultListenAddrPort = "127.0.0.1:8080"
	ascii_icon            = `

                         ____==========_______
              _--____   |    | ""  " "|       \
             /  )8}  ^^^| 0  |  =     |  o  0  |
           </_ +-==B vvv|""  |  =     | '  "" "|
              \_____/   |____|________|________|
                       (_(  )\________/___(  )__)
                         |\  \            /  /\    
                         | \  \          /  /\ \      
                         | |\  \        /  /  \ \    
                         (  )(  )       (  \   (  )
                          \  / /        \  \   \  \    
                           \|  |\        \  \  |  |
                            |  | )____    \  \ \  )___      
                            (  )  /  /    (  )  (/  /
                           /___\ /__/     /___\ /__/
=====================================================
                                The AT-AT, By Core21

`
)

// Packaged all settings
type Config struct {
	global  Global
	server  Server
	db      Database
	cache   Cache
}

type Global struct {
	debug bool
	cacheEable bool
	logFile string
}

// Packaged all database settings
type Database struct {
	addr     string
	protocol string
	user     string
	password string
	dbname   string
}

// Packaged all Server settings
type Server struct {
	listenAddr    string
	port          string
	cpuCore       int
	gzip          bool
	jsonIndent    bool
	statusService bool
}

type Cache struct {
	addr string
	protocol string
}

// Define a global config varible
var config Config

// Parse the config.conf file
func parseConfig(configPath string) {
	var err error
	var exists bool
	exists, err = IsFileExists(configPath)

	if err != nil && !exists {
		panic("Fail to load configuration file: " + err.Error())
	}
	var cfg *goconfig.ConfigFile
	cfg, err = goconfig.LoadConfigFile(configPath)
	if err != nil {
		panic("Fail to load configuration file: " + err.Error())
	}

	// Parse the global setcion
	if cfg.MustValue("global", "debug") == "true" {
		config.global.debug = true
	} else {
		config.global.debug = false
	}
	if cfg.MustValue("global", "cache_enable") == "true" {
		config.global.cacheEable = true
	} else {
		config.global.cacheEable = false
	}
	config.global.logFile = cfg.MustValue("global", "logfile")


	// Parse the server setcion
	if cfg.MustValue("server", "gzip") == "true" {
		config.server.gzip = true
	} else {
		config.server.gzip = false
	}

	if cfg.MustValue("server", "status_service") == "true" {
		config.server.statusService = true
	} else {
		config.server.statusService = false
	}

	if cfg.MustValue("server", "json_indent") == "true" {
		config.server.jsonIndent = false
	} else {
		config.server.jsonIndent = true
	}
	config.server.cpuCore, err = strconv.Atoi(cfg.MustValue("server", "core"))
	checkErr(err)
	config.server.listenAddr = cfg.MustValue("server", "listen")
	config.server.port = cfg.MustValue("server", "port")


	// Parse the database setcion
	config.db.addr = cfg.MustValue("database", "addr")
	config.db.protocol = cfg.MustValue("database", "protocol")
	config.db.user = cfg.MustValue("database", "user")
	config.db.password = cfg.MustValue("database", "password")
	config.db.dbname = cfg.MustValue("database", "dbname")


	// Parse the cache section
	config.cache.addr = cfg.MustValue("cache", "addr")
	config.cache.protocol = cfg.MustValue("cache", "protocol")

	// fmt.Println("cpucore: ",config.Server.CPUCore)
	// fmt.Println("port: ",config.Server.Port)
	// fmt.Println("Status: ", config.Server.StatusService)
	// fmt.Println("Gzip: ", config.Server.Gzip)
	// log.Println("debug: ", config.Server.JsonIndent)
}
