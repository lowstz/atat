package main

import (
	"os"
	"fmt"
	"io"
	"log"
)

// Print the atat server version.
func printVersion() {
	fmt.Println("AT-AT version", version)
}

// Print a super AT-AT ASCII ICON.
// The detail about AT-AT you can found from wikipedia:
// http://en.wikipedia.org/wiki/Walker_(Star_Wars)#All_Terrain_Armored_Transport_.28AT-AT.29
func printAsciiIcon() {
	fmt.Println(ascii_icon)
}

// return custom log.Logger
func getLogFile(logpath string) (io.Writer) {
	logFile := os.Stderr
	if logpath != "" {
		if f, err := os.OpenFile(logpath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600); err != nil {
			fmt.Printf("Can't open log file, logging to stderr: %v\n", err)
		}else {
			logFile = f
		}
	}
	return logFile
}

func serverStartInfoLog(logger io.Writer) (*log.Logger) {
	startInfoLogger := log.New(logger, "[Start] ", log.LstdFlags)
	return startInfoLogger
}



















