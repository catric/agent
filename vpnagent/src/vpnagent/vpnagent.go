package main

import (
	_ "daemon"
	"fmt"
	"logging"
	"net/http"
	"net/url"
	"os"
	"time"
	"io"
	"runtime"
	"encoding/json"
	"runtime/debug"
)

const (

	SERVER_NAME = "vpnagent"
	SERVER_VERSION = "1.00n"
	SERVER_PORT = "62688"
)

// 初始化日志
var Log *logging.Logger = nil

func CreateLogger(loggerName string, filePrefix string, logLv string) *logging.Logger {
	Log = logging.MustGetLogger(loggerName)
	Log.CustomSettings("[%{time:2006-01-02 15:04:05.000}][%{shortfunc}][%{level:.4s}] %{message}", true, "log/" + filePrefix, logLv)
	return Log
}

func initLog() {
	Log = CreateLogger(SERVER_NAME, SERVER_NAME, "DEBUG")
}

func freeLog() {
}

func CatchException() {
	if err := recover(); err != nil {
		Log.Critical("err = %v, stack = %v", err , string(debug.Stack()))
	}
}

func SerializeToJSON(st interface{}) string {
	ba, _ := json.Marshal(st)
	jsonstr := string(ba)
	
	return jsonstr
}

func UnserializeFromJSON(jsonstr string, st interface{}) {
     	json.Unmarshal([]byte(jsonstr), st)
}

func ParamEncode(data string) string {
	return url.QueryEscape(data)
}

func ParamDecode(data string) (string, error) {

	sdata, err := url.QueryUnescape(data)
	if err != nil {
		Log.Error("url.QueryUnescape err = %v", err)
		return "", err
	}
	
	return sdata, nil
}

func UseMaxCpu() {
	// multiple cups using
	runtime.GOMAXPROCS(runtime.NumCPU())
	
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
//	if r.URL.Path == "/" {
//		http.Redirect(w, r, "/agent", http.StatusFound)
//		return
//	}
	html := `Not Found 404 404 404 404`
	io.WriteString(w, html)
}

func AgentHandler(w http.ResponseWriter, r *http.Request) {

//	if r.URL.Path == "/" {
//		http.Redirect(w, r, "/agent", http.StatusFound)
//		return
//	}

	

	html := `Not Found 404 404 404 404`
	io.WriteString(w, html)
}

func runServer() {

	//set request handler
	http.HandleFunc("/agent", AgentHandler)
	http.HandleFunc("/", NotFoundHandler)

	err := http.ListenAndServe(":" + SERVER_PORT, nil)
	if err != nil {
		Log.Panic("ListenAndServe: " + err.Error())
	}
}

func startup() {

	Log.Info(SERVER_NAME + " System starting ... ...")

	Log.Info("HTTP Server starting ... ...")
	Log.Info("	PORT = " + SERVER_PORT)
	go runServer()
	
}

func shutdown() {
	Log.Info("System shutdown ... ...")
}


var Shutdown bool = false

func main() {
	// all-in-one mode, distributed-mode
	// daemon-mode, front-mode
	initLog()

	defer CatchException()

	Log.Info(SERVER_NAME + " Ver=" + SERVER_VERSION)
	startup()
	
	// dead loop for console input
	argc := len(os.Args)
	if argc == 1 { // console mode
		for {
			cmd := ""

			fmt.Println("Console : ")
			fmt.Scanln(&cmd)

			if cmd == "exit" {
				shutdown()
				os.Exit(1)
			} else {
				fmt.Println(cmd)
			}
		}
	} else { // daemon mode
		for {
			if Shutdown {
				shutdown()
				os.Exit(1)
			}

			time.Sleep(2 * time.Second)
		}
	}
}
