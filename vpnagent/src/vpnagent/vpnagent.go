package main

import (
	_ "daemon"
	"fmt"
	"logging"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"
	"io"
//	"io/ioutil"
	"runtime"
	"path/filepath"
	"encoding/json"
	"runtime/debug"
	"github.com/Unknwon/goconfig"
)

const (

	SERVER_NAME = "vpnagent"
	SERVER_VERSION = "1.00n"
	//SERVER_PORT = "62688"
	SERVER_INI = "agent.ini"
)

const (
	METHOD_ADDUSER = "adduser"
	METHOD_BANUSER = "banuser"
	METHOD_DELUSER = "deluser"
	METHOD_UPDUSER = "upduser"
	METHOD_SYNUSER = "synuser"
)

// 配置文件内存结构
type TAgentConfig struct {
    Host        string
    Port        string
    Loglevel    string
    ManagePassword string
	CommKey		string
}

type TUser struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type TUserList struct {
	Users []TUser `json:"Users"`
}

type TRetResp struct {
	Retcode string `json:"Retcode"`
	Retdesc string `json:"Retdesc"`
}

type TNodeList struct {
	Nodes []TNode `json:"Nodes"`
}

type TNode struct {
	Ip string `json:"Ip"`
	Port string `json:"Port"`
}

var VAConfig TAgentConfig
var CurrPath string = ""

// 初始化日志
var Log *logging.Logger = nil

func CreateLogger(loggerName string, filePrefix string, logLv string) *logging.Logger {
	Log = logging.MustGetLogger(loggerName)
	Log.CustomSettings("[%{time:2006-01-02 15:04:05.000}][%{shortfunc}][%{level:.4s}] %{message}", true, "log/" + filePrefix, logLv)
	return Log
}

func initLog() {
	Log = CreateLogger(SERVER_NAME, SERVER_NAME, VAConfig.Loglevel)
}

func freeLog() {
}

func CatchException() {
	if err := recover(); err != nil {
		Log.Critical("err = %v, stack = %v", err , string(debug.Stack()))
	}
}


func loadConfig() {
    conf, err := goconfig.LoadConfigFile("config/" + SERVER_INI)
    if err == nil {
        if conf != nil {
            //runMode, _ := conf.GetValue("main", "run_mode")
            VAConfig.Host, _ = conf.GetValue("main", "host")
            VAConfig.Port, _ = conf.GetValue("main", "port")
            VAConfig.Loglevel, _ = conf.GetValue("main", "log_level")
            VAConfig.CommKey, _ = conf.GetValue("main", "comm_key")
            VAConfig.ManagePassword,_ = conf.GetValue("main", "manage_password")
        } else {
            fmt.Println("Load config file failed! ")
            os.Exit(-1)
        }
    } else {
        fmt.Println("Load config file failed! please check the config path!")
        os.Exit(-1)
    }
}


func allowCrossRequest(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
    w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
    w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS");
}

func contentTypeJSON(w http.ResponseWriter) {
    w.Header().Set("Content-Type", "application/json") //返回数据格式是json
}

func SetHeaderAllowCrossRequest(w http.ResponseWriter) {
    allowCrossRequest(w)
    contentTypeJSON(w)
}

func SerializeToJSON(st interface{}) string {
	ba, _ := json.Marshal(st)
	jsonstr := string(ba)
	
	return jsonstr
}

func UnserializeFromJSON(jsonstr string, st interface{}) {
     	json.Unmarshal([]byte(jsonstr), st)
}

func MakeResp(retcode string) string {
    var resp TRetResp
    resp.Retcode = retcode
    resp.Retdesc = ErrMsg[retcode]

    baresp, _ := json.Marshal(resp)
    strresp := string(baresp)

    return strresp
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

func ParseParams(data string) map[string]string {

    m := make(map[string]string)
    
	vals, err := url.ParseQuery(data)
	if err != nil {
		Log.Error("Parse params error = %v", err)
		return nil
	}
        
	for k,v := range vals {
		m[k] = v[0]
	}
    
    Log.Debug("Parse params m = %v", m)

    return m
}

// Entry parameters check
func CheckParams(data string) (retcode string, paramsMap map[string]string) {

    paras := data
    
    params := ParseParams(string(paras))
    method := params["m"]
    
    if method == METHOD_ADDUSER {

    } else if method == METHOD_BANUSER {
    
    } else if method == METHOD_UPDUSER {

    } else if method == METHOD_SYNUSER {

    } else if method == METHOD_DELUSER {

    } else {
        return RC_INVALID_REQUEST, nil
    }
    
    return RC_SUC, params
}

func makeBatchFile(batchCmd string) (bool, string) {
	
	Log.Debug("makeBatchFile Entering... ...")
	
	//filename := "./cmdfile/" + RandomAlphabetic(8) + ".txt"
	filename := CurrPath + "/cmdfile/" + RandomAlphabetic(8) + ".txt"
	Log.Debug("filename = %v", filename)
	
	//file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	file, err := os.Create(filename)
	if err != nil {
		Log.Error("err=%v", err)
		return false, ""
	}
	
	Log.Debug("create batch cmd file... ...file=%v", file)
	
	file.WriteString(batchCmd)
	file.Sync()
	file.Close()
	
	return true, filename
}

func runCmd(cmdstr string) bool {
	
	cmd := exec.Command(cmdstr)
	err := cmd.Run()
	if err != nil {
		Log.Error("err=%v", err)
		return false
	}
	
	return true
}

func executeBatchCommands(nodelist TNodeList, batchCmd string) string {

	ok, filename := makeBatchFile(batchCmd)
	if !ok {
		Log.Error("make batch cmd file failed!")
		return RC_MAKE_BATCH_CMD_FAILED
	}
	
	defer func() {

		Log.Info("execute command over ... ...")

//		os.Remove(filename)
	}()

	ns := nodelist.Nodes
	for idx := range ns {
		node := ns[idx]		

		cmdstr := CurrPath + "/cmd/vpncmd " + node.Ip + ":443" + " /SERVER /PASSWORD:" + VAConfig.ManagePassword + " /IN:" + filename
		Log.Debug("cmdstr = %v", cmdstr)
		
		//go runCmd(cmdstr)
		//go func() bool {
		func() bool {
			cmd := exec.Command(CurrPath + "/cmd/vpncmd", node.Ip + ":443", "/SERVER", "/PASSWORD:" + VAConfig.ManagePassword, " /IN:" + filename)
			output, err := cmd.Output()
			if err != nil {
				Log.Error("err=%v", err)
				return false
			}
			
			Log.Info("command output = %v", string(output))
			
			return true
		}()
		
	}
	
	return RC_SUC
}

func makeAdduserCommands(uname string, pwd string) []string {

	var cmds []string
	
	cmd := " UserCreate "+ uname +" /GROUP:none /REALNAME:" + uname + " /NOTE:none "
	cmds = append(cmds, cmd)
	
	cmd =  " UserPasswordSet "+uname +" /PASSWORD:"+ pwd 
	cmds = append(cmds, cmd)
	
	return cmds
}

func makeBanuserCommands(uname string) []string {

	var cmds []string
	
	cmd := " UserPolicySet "+ uname +" /NAME:Access /VALUE:no "
	cmds = append(cmds, cmd)
	
	return cmds
}

func makeDeluserCommands(uname string) []string {

	var cmds []string
	
	cmd := " UserDel "+ uname
	cmds = append(cmds, cmd)
	
	return cmds
}

func makeUpduserCommands(uname string, pwd string, valid string) []string {

	var cmds []string
	
	cmd := " UserList"
	cmds = append(cmds, cmd)
	
	if len(pwd) > 0 {
		cmd := " UserPasswordSet "+uname +" /PASSWORD:"+ pwd
		cmds = append(cmds, cmd)
	}
	
	if len(valid) > 0 {
		cmd := " UserExpiresSet " + uname + " /EXPIRES:" + valid
		cmds = append(cmds, cmd)
	} 
	
	return cmds
}

func executeCommands(nodelist TNodeList, cmds []string) string {
	
	defer func() {

		Log.Info("execute command over ... ...")

//		os.Remove(filename)
	}()

	ns := nodelist.Nodes
	for idx := range ns {
		node := ns[idx]

		go func() bool {
			for j := range cmds {
				cmdstr := CurrPath + "/cmd/vpncmd " + node.Ip + ":443" + " /SERVER /PASSWORD:" + VAConfig.ManagePassword + " /ADMINHUB:vpn /CMD " + cmds[j]
				Log.Debug("cmdstr = %v", cmdstr)
				
				cmd := exec.Command(CurrPath + "/cmd/vpncmd", node.Ip + ":443", "/SERVER", "/PASSWORD:" + VAConfig.ManagePassword, "/ADMINHUB:vpn", "/CMD " + cmds[j])
				output, err := cmd.Output()
				if err != nil {
					Log.Error("err=%v", err)
					return false
				}
				
				Log.Info("command output = %v", string(output))
			}
			
			return true
		}()
	}
	
	return RC_SUC
}

func AddUserHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
    
    Log.Info("Entering AddUserHandler ...")
    
    nodes := params["nodes"]
    uname := params["uname"]
//    utype := params["utype"]
    pwd   := params["pwd"]
//    valid := params["valid"]

	var nodelist TNodeList
	UnserializeFromJSON(nodes, &nodelist)
	
//	batchCmd := "\n hub vpn \n" + 
//			" UserCreate "+ uname +" /GROUP:none /REALNAME:" + uname + " /NOTE:none \n" +
//			" UserPasswordSet "+uname +" /PASSWORD:"+ pwd +" \n"
	
	cmds := makeAdduserCommands(uname, pwd)
	Log.Debug("adduser cmds = %v  ", cmds)
	
	rc := executeCommands(nodelist, cmds)
	jsonresp := MakeResp(rc)
	io.WriteString(w, jsonresp)
}

func BanUserHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
    
    Log.Info("Entering BanUserHandler ...")
    
    nodes := params["nodes"]
    uname := params["uname"]

	var nodelist TNodeList
	UnserializeFromJSON(nodes, &nodelist)
	
	cmds := makeBanuserCommands(uname)
	Log.Debug("banuser cmds = %v  ", cmds)
	
	rc := executeCommands(nodelist, cmds)
	jsonresp := MakeResp(rc)
	io.WriteString(w, jsonresp)
}

func DelUserHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
    
    Log.Info("Entering DelUserHandler ...")
    
    nodes := params["nodes"]
    uname := params["uname"]

	var nodelist TNodeList
	UnserializeFromJSON(nodes, &nodelist)
	
	cmds := makeDeluserCommands(uname)
	Log.Debug("banuser cmds = %v  ", cmds)
	
	rc := executeCommands(nodelist, cmds)
	jsonresp := MakeResp(rc)
	io.WriteString(w, jsonresp)
}

func UpdUserHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
    
    Log.Info("Entering UpdUserHandler ...")
    
    nodes := params["nodes"]
    uname := params["uname"]
//    utype := params["utype"]
    pwd   := params["pwd"]
    valid := params["valid"]

	var nodelist TNodeList
	UnserializeFromJSON(nodes, &nodelist)
	
	cmds := makeUpduserCommands(uname, pwd, valid)
	Log.Debug("upduser cmds = %v  ", cmds)
	
	rc := executeCommands(nodelist, cmds)
	jsonresp := MakeResp(rc)
	io.WriteString(w, jsonresp)
}

func SynUserHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
    
    Log.Info("Entering SynUserHandler ...")
    
//    node := params["node"]
//    users := params["users"]
//    
//    var userlist TUserList
//    UnserializeFromJSON(users, &userlist)
//    
//	batchCmd := "\n hub vpn \n"
//    
//    us := userlist.Users
//    for idx := range us {
//    	u := us[idx]
//    	
//    	batchCmd = batchCmd + 
//    				" UserCreate "+ u.Username +" \n" +
//					" UserPasswordSet "+ u.Username +" /PASSWORD:" + u.Password +" \n" 
//    	
//    }
//    
//    Log.Debug("batch cmd = %v", batchCmd)
//    
//    var n TNode
//   	var nodelist TNodeList
//   	n.Ip = node
//   	nodelist.Nodes = append(nodelist.Nodes, n)
//   	
//	rc := executeCommands(nodelist, batchCmd)
//	jsonresp := MakeResp(rc)
//	io.WriteString(w, jsonresp)
}


func AgentHandler(w http.ResponseWriter, r *http.Request) {

    Log.Info("Entering AgentHandler ...")
    SetHeaderAllowCrossRequest(w)
    
    defer CatchException()

    err := r.ParseForm()
    if err != nil {
        Log.Error("err = " + err.Error())
        strresp := MakeResp(RC_PARAM_ERROR)
        io.WriteString(w, strresp)
        return
    }
    
    data := r.FormValue("data")
    Log.Debug("data = " + data)
    
    // 预先检查参数
    retcode, params := CheckParams(data)
    if !IsSuccess(retcode) {
        Log.Error("err = " + ErrMsg[retcode])
        strresp := MakeResp(retcode)
        io.WriteString(w, strresp)
    }
    
    method := params["m"]

    if method == METHOD_ADDUSER {
		AddUserHandler(w, r, params)
    } else if method == METHOD_BANUSER {
    	BanUserHandler(w, r, params)
    } else if method == METHOD_UPDUSER {
		UpdUserHandler(w, r, params)
    } else if method == METHOD_SYNUSER {
		SynUserHandler(w, r, params)
    } else if method == METHOD_DELUSER {
		DelUserHandler(w, r, params)
    } else {
        Log.Error("err = " + ErrMsg[RC_INVALID_REQUEST])
        strresp := MakeResp(RC_INVALID_REQUEST)
        io.WriteString(w, strresp)
    }
    
    return
}

func RunServer() {

	//set request handler
	http.HandleFunc("/agent", AgentHandler)
	http.HandleFunc("/", NotFoundHandler)

	err := http.ListenAndServe(":" + VAConfig.Port, nil)
	if err != nil {
		Log.Panic("ListenAndServe: " + err.Error())
	}
}

func startup() {

	Log.Info(SERVER_NAME + " System starting ... ...")

	Log.Info("HTTP Server starting ... ...")
	Log.Info("	PORT = " + VAConfig.Port)
	go RunServer()
	
}

func GetCurrPath() string {
//	file, _ := exec.LookPath(os.Args[0])
//    path, _ := filepath.Abs(file)

    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))  
    if err != nil {  
        Log.Error("err=%v", err)
        
        return ""
    }  

    return dir
}

func shutdown() {
	Log.Info("System shutdown ... ...")
}

var Shutdown bool = false

func main() {
	// all-in-one mode, distributed-mode
	// daemon-mode, front-mode
	loadConfig()
	initLog()
	
	defer CatchException()
	
	Log.Info("########################")
	Log.Info(SERVER_NAME + " Ver=" + SERVER_VERSION)
	Log.Info("########################")

	CurrPath = GetCurrPath()
	Log.Info("Current Path = %v", CurrPath)
	
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
