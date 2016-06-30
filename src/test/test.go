package  main 

import (
        "net"
        "os"
        "io"
        "fmt"
        "time"
        "encoding/json"
        "math/rand"
)

import (
        rpc "rpc_framework"
)

type Msg struct{
    Method  string  `json:"method"`
    Topic   string  `json:"topic"`
    Data    string  `json:"data"`
    Ts      int64   `json:"ts"`
    Sign    string  `json:"sign"`
    Expired int64   `json:"expired"`
}
func main() {
    conn, err := net.Dial("tcp", "127.0.0.1:8886")
    defer conn.Close()
    checkError(err)
    
    nshead := rpc.NsHead{}
    
    var i int = 0 
    var NUM int = 100000
    startTime := time.Now()
     
    for {   
        i ++
        nshead.LogId = uint32(rand.Intn(NUM)) 

        msg := &Msg{}
        msg.Method = "producer.push"
        msg.Topic = "topic_test"
        msg.Data = fmt.Sprintf("test.client.payload:%d", i)
        msg.Ts = time.Now().Unix()
        msg.Sign = "sign_test"
        msg.Expired = 0 
        
        payload, err:= json.Marshal(msg)
        data, err := nshead.Marshal([]byte(payload))
        checkError(err)
        
        _, err = conn.Write(data)
        checkError(err)
    
        resp, err := readFully(conn)
        checkError(err)
        
        if i % NUM == 0 {
            fmt.Printf("request:%d, time_cost:%v\n", NUM, time.Now().Sub(startTime))
            fmt.Println("response: ", string(resp))
            startTime = time.Now()
            
            i = 1 
        }
    }
    
    os.Exit(0)
}

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}

func readFully(conn net.Conn) ([]byte, error) {

    nh := rpc.NsHead{}
    headLen, _ := nh.GetHeaderLen()
    headBuf := make([]byte, headLen)
    _, err := io.ReadFull(conn, headBuf)
    checkError(err)

    //_, _ := nh.GetLogId(headBuf)
    bodyLen, err := nh.GetBodyLen(headBuf)
    checkError(err)
    
    bodyBuf := make([]byte, bodyLen)
    _, err = io.ReadFull(conn, bodyBuf)
    checkError(err)
    
    return  bodyBuf, nil
}
