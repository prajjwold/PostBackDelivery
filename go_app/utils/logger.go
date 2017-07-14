package utils

import (
	"go/build"
	"io"
	"log"
	"os"
)

const INFO_LOG_FILE = "delivery_agent.log"
const POSTBACK_LOG_FILE = "postback_request_response.log"

var (
	//Logger specifically for postback request response trace
	PostBack *log.Logger
	//Logger for general trace and error logs
	Info *log.Logger
)

func SetupLogger() {
	gopath := build.Default.GOPATH
	logdir := gopath + "/src/postbackdelivery/logs"
	mode := os.ModePerm
	if _, err := os.Stat(logdir); os.IsNotExist(err) {
		os.Mkdir(logdir, mode)
	}

	logpath := logdir + "/" + INFO_LOG_FILE
	file, err := os.Create(logpath)
	if err != nil {
		log.Fatalln("Failed to create log file", INFO_LOG_FILE, ":", err)
	}
	multi := io.MultiWriter(file, os.Stdout)
	Info = log.New(multi,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	logpath = logdir + "/" + POSTBACK_LOG_FILE
	file, err = os.Create(logpath)
	if err != nil {
		log.Fatalln("Failed to create log file", POSTBACK_LOG_FILE, ":", err)
	}
	multi = io.MultiWriter(file, os.Stdout)
	PostBack = log.New(multi,
		"POSTBACK: ",
		log.Ldate|log.Ltime|log.Lshortfile)

}
