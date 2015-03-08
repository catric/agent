package main

import (

)


type TErrMsg struct {
    code string
    desc string
}

type ErrCode string

const (
    RC_SUC = "0"
    
    RC_PARAM_ERROR = "10001"
    RC_MAKE_BATCH_CMD_FAILED = "10002"
    
    RC_GET_HASHMAP_ERROR = "10003"
    RC_USER_CANNOT_BE_EMPTY = "10004"
    RC_FILE_TYPE_NOT_ALLOWED = "10005"
    RC_FILE_UPLOAD_FAILED = "10006"
    RC_FILE_TOO_LARGE      = "10007"
    RC_INVALID_REQUEST    = "10008"
    RC_UNKNOWN_REPORT_TYPE = "10009"
    RC_UNKNOWN_NODE_TYPE = "10010"
    RC_NODE_NAME_NOT_BE_EMPTY = "10011"
    RC_NODE_REPORT_FAILED = "10012"
    RC_SEND_TO_GATE_FAILED = "10013"
    RC_SEND_MSG_FAILED = "10014"
    RC_ADD_MEMBER_FAILED = "10015"
    RC_USER_OFFLINE = "10016"
    RC_USER_NOT_EXISTS = "10017"
    RC_DB_ERROR        = "10018"
    RC_CS_NOT_EXISTS = "10019"
    RC_CS_NAME_OR_PASSWORD_MISMATCH = "10020"
    RC_GET_USER_GROUPS_FAILED = "10021"
    RC_NO_CS_AVALIABLE = "10022"
    RC_PROMO_CODE_NOT_EXISTS = "10023"
    RC_SESSION_ID_MISMATCH = "10024"
    RC_NO_MANAGER_PRIVILEGE = "10025"
    RC_CS_NOT_AVAILABLE = "10026"
    RC_USER_NOT_IN_CACHE = "10027"
    RC_UNKNOWN_MONITOR_TYPE = "10028"
    RC_MOV_MEMBER_FAILED = "10029"
    RC_GROUP_NOT_EXISTS = "10030"
    RC_NOT_CROWD_MEMBER = "10031"
    
    // 20000 + err code for NiFileServer
    
    
    
    // 30000 + err code for NiGateServer
    RC_NIDECODE_FAILED = "30001"
    RC_CHAT_SERVER_ERROR = "30002"
    RC_UNKNOWN_RECV_TYPE = "30003"
    RC_PARAMS_CANNOT_BE_EMPTY = "30004"
)

var ErrMsg = map[string]string {
    RC_SUC : "SUCCESS",
    RC_PARAM_ERROR : "Parameter error",
    RC_MAKE_BATCH_CMD_FAILED : "Make batch cmd file failed",
    
    
    
    RC_USER_CANNOT_BE_EMPTY : "User can not be empty",
    RC_FILE_TYPE_NOT_ALLOWED : "This type of file not be allowed to upload",
    RC_FILE_UPLOAD_FAILED : "File upload failed",
    RC_FILE_TOO_LARGE : "File too large",
    RC_INVALID_REQUEST : "Invalid request",
    RC_UNKNOWN_REPORT_TYPE : "Unknown report type",
    RC_UNKNOWN_NODE_TYPE : "Unknown node type",
    RC_NODE_REPORT_FAILED : "Node report to charserver failed",
    RC_SEND_TO_GATE_FAILED : "Send to gate failed ",
    RC_SEND_MSG_FAILED : "Send message failed",
    RC_ADD_MEMBER_FAILED : "Add member to group failed",
    RC_USER_OFFLINE : "User offline",
    RC_USER_NOT_EXISTS : "User does not exists",
    RC_DB_ERROR : "Db error",
    RC_CS_NOT_EXISTS : "Cs does not exists",
    RC_CS_NAME_OR_PASSWORD_MISMATCH : "Name or password mismatch",
    RC_GET_USER_GROUPS_FAILED : "Get user groups failed",
    RC_NO_CS_AVALIABLE : "No cs available",
    RC_PROMO_CODE_NOT_EXISTS : "Promo code not exists",
    RC_SESSION_ID_MISMATCH : "Session ID mismatch",
    RC_NO_MANAGER_PRIVILEGE : "No manager privilege",
    RC_CS_NOT_AVAILABLE : "CS not available",
    RC_USER_NOT_IN_CACHE : "User not in cache",
    RC_UNKNOWN_MONITOR_TYPE : "Unknown monitor type",
    RC_MOV_MEMBER_FAILED : "Move member to group failed",
    RC_GROUP_NOT_EXISTS : "Group not exists",
    RC_NOT_CROWD_MEMBER : "Not crowd member",
    
    // 20000 + err code for NiFileServer


    // 30000 + err code for NiGateServer
    RC_NIDECODE_FAILED : "NiDecode data failed",
    RC_CHAT_SERVER_ERROR    : "Chat server error",
    RC_UNKNOWN_RECV_TYPE : "Unknown recv type",
    RC_PARAMS_CANNOT_BE_EMPTY : "Params cannot be empty",
}

func IsSuccess(code string) bool {
    if code == RC_SUC {
        return true
    } else {
        return false
    }
}




