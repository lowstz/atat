package main

import (
	"fmt"
	"github.com/Unknwon/goconfig"
)

const (
	version           = "0.1"
	defaultListenAddrPort = "127.0.0.1:8080"
)

// Packaged all settings
type Config struct {
	Debug  bool
	Server Server
	Db     Database
}

// Packaged all database settings
type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

// Packaged all Server settings
type Server struct {
	ListenAddr string
	Port       string
}

// Define a global config varible
var config Config

// Print the atat server version.
func printVersion() {
	fmt.Println("AT-AT version", version)
}

// Print a super AT-AT ASCII ICON.
// The detail about AT-AT you can found from wikipedia:
// http://en.wikipedia.org/wiki/Walker_(Star_Wars)#All_Terrain_Armored_Transport_.28AT-AT.29
func printAsciiIcon() {
	fmt.Println(`
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
`)
}

// Parse the config.conf file
func parseConfig(configPath string) {
	var err error
	var exists bool
	exists, err = isFileExists(configPath)

	if err != nil && !exists {
		panic("Fail to load configuration file: " + err.Error())
	}
	var cfg *goconfig.ConfigFile
	cfg, err = goconfig.LoadConfigFile(configPath)
	if err != nil {
		panic("Fail to load configuration file: " + err.Error())
	}

	if cfg.MustValue("global", "debug") == "true" {
		config.Debug = true
	}
	config.Server.ListenAddr = cfg.MustValue("server", "listen")
	config.Server.Port = cfg.MustValue("server", "port")
	config.Db.Host = cfg.MustValue("database", "host")
	config.Db.Port = cfg.MustValue("database", "port")
	config.Db.User = cfg.MustValue("database", "user")
	config.Db.Password = cfg.MustValue("database", "password")
	config.Db.Dbname = cfg.MustValue("database", "dbname")
}
