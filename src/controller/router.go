package controller

import (
    "encoding/json"
    "strings"
)

import(
    "utils"
    "action"
)

type Router struct {}

type Methods struct {
    Method  string  `json:"method,omitempty"` 
}

func (handle *Router) Route(logid uint32, input []byte) (resp []byte, err error) {
    defer func() {
        if e, ok := recover().(error); ok {
            resp = []byte(e.Error())
            return
        }
    }() 
    
    logger := utils.Global.Logger
    logbuf := utils.NewLogBuffer(logid)
    logbuf.WriteLog(" [request:%s]", string(input)) 
   
    var method Methods
    var output string
    err = json.Unmarshal(input, &method)
    if err != nil {
       return nil, err 
    } 
    
    subMethods := strings.Split(method.Method, ".")
    if len(subMethods) < 2 {
        panic(&utils.SysError{logbuf, "err.method_not_support"})
    }
     
    act := subMethods[0]
    subAct := subMethods[1]
    switch  act{
        case "producer":
            output = action.Producer(subAct, logbuf, input)
        default:
            panic(&utils.SysError{logbuf, "err.method_not_support"})
    }
    logbuf.WriteLog(" [error_code:0]")
    logger.UbLogNotice(logbuf.String())    
    return []byte(output), nil
}
