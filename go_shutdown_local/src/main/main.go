package main

import(
    _"syscall"
    "fmt"
    "os/exec"
    "log"
    "io"
    "net"
    _"os"
    "bytes"
    "net/http"
    "io/ioutil"
)

// TCP Server


func main() {
    // 1) Send the endpoint INFO to server
    reqBody := bytes.NewBufferString(GetLocalIP())
    fmt.Println(GetLocalIP())

    // 34.225.204.24
    recv, err := http.Post("http://127.0.0.1:8004/api/v1/notebook/turnon", "text/plain", reqBody) 
    fmt.Println("Send turnon signal to server")
    if err != nil{
        
        fmt.Println(err.Error())
    }

    fmt.Println(recv)
    fmt.Println("recieve Response")
    
    var url string = "http://127.0.0.1:8004/api/v1/long?timeout=3600&category=cmd"
    resp, err := http.Get(url)
    if err != nil {
        panic(err)
    }
 
    defer resp.Body.Close()
 
    // 결과 출력
    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
    // {"events":[{"timestamp":1565616607427,"category":"farmss","data":"shutdown plz"}]}
    // {"timeout":"no events before timeout","timestamp":1565615537077}


    fmt.Printf("%s\n", string(data))


}

func ConnHandler(conn net.Conn) {
    recvBuf := make([]byte, 4096) // receive buffer: 4kB
    for {
        n, err := conn.Read(recvBuf)
        if nil != err {
            if io.EOF == err {
                log.Printf("connection is closed from client; %v", conn.RemoteAddr().String())
                return
            }
            log.Printf("fail to receive data; err: %v", err)
            return
        }

        if 0 < n {
            data := recvBuf[:n]
            log.Println(string(data))
            
            if string(data) == "shutdown"{
                Shutdown()
            }

        }
    }
}

func GetLocalIP() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return ""
    }
    for _, address := range addrs {
        // check the address type and if it is not a loopback the display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return ""
}

// Shutdown the this computer
func Shutdown(){
    
    cmd := exec.Command("shutdown", "now")
    err := cmd.Start()
    if err != nil {
        log.Fatal(err)
    }
}

