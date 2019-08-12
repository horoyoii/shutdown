// Server

/*
 Created by Horoyoii on 2019.07.29
*/


/*
 1. receive the shutdown command from android
 2. receive the endpoint info from local program

 3. send notebook status info to android
 4. send the shutdown command to local program

*/
package main

import (
    "log"
    _"net"
    "github.com/gin-gonic/gin"
    "fmt"
    "net/http"
    "github.com/jcuga/golongpoll"
    
)


var CONN_IP_NOTEBOOK string
var STATUS_NOTEBOOK bool
var flag bool = false


func main() {
    router := gin.Default()
    
    manager, err := golongpoll.StartLongpoll(golongpoll.Options{
        LoggingEnabled: true,
        MaxLongpollTimeoutSeconds: 3600,
        })

    if err != nil {
        log.Fatalf("Failed to create manager: %q", err)
    }


    android := router.Group("api/v1/android")
    {
        android.GET("status", Status)
        android.GET("/shutdown", func(c *gin.Context) {


            fmt.Println("shutdonw commmand comes in from android")

            // If notebook is not running...
            if STATUS_NOTEBOOK != true {
                fmt.Println("notebook is not running...")
                return
            }

            fmt.Println("notebook is running... and send event response")
            
            STATUS_NOTEBOOK = false

            manager.Publish("cmd", "shutdown plz")

            c.JSON(200, gin.H{
			    "message": "pong",
		    })
	    })
    }


    notebook := router.Group("api/v1/notebook")
    {
        notebook.POST("/turnon", TurnOnEndpoint)
        notebook.GET("/ping", PingToNotebook)
    }



    router.GET("/api/v1/ping", PingOfServer)
    router.GET("/api/v1/long", gin.WrapF(manager.SubscriptionHandler))
    
    
    router.Run(":8004")
}

func generateRandomEvents(lpManager *golongpoll.LongpollManager) {
    for {
        if flag {
            lpManager.Publish("farmss", "shutdown plz")
        }
    }
}


func PingToNotebook(c *gin.Context){
    fmt.Println("Sent to notebook and waiting... response") 
    resp, err := http.Get("http://"+CONN_IP_NOTEBOOK+":8003/ping")
    if err != nil{
        fmt.Println(err.Error())
    }
    defer resp.Body.Close()
    fmt.Println(resp.Body) 
    c.String(http.StatusOK, "OK ") 
}


func PingOfServer(c *gin.Context){



    c.JSON(200, gin.H{
        "message":"pong",
        "notebook_state":STATUS_NOTEBOOK,
        "notebook_ip":CONN_IP_NOTEBOOK,
    })

}


/* Send status of notebook to android 
*/
func Status(c *gin.Context){
    var message string

    if STATUS_NOTEBOOK {
        message = "on"
    }else{
        message = "off"
    }

    c.String(http.StatusOK, message)
}



/* send shutdown command to notebook
*/
func Shutdown(c *gin.Context){
    fmt.Println("shutdonw commmand comes in from android")

    // If notebook is not running...
    if STATUS_NOTEBOOK != true {
        fmt.Println("notebook is not running...")
        return
    }

    fmt.Println("notebook is running... and send event response")
    flag = true

    STATUS_NOTEBOOK = false
}



/* Get the IP of notebook from notebook and Save this info.
*/
func TurnOnEndpoint(c *gin.Context){

    // notebook IP is sent via POST body
    b, err := c.GetRawData()
    if err != nil{
        panic(err)
    }
    CONN_IP_NOTEBOOK = string(b)
    fmt.Println(string(b))

    // Set Status ON
    STATUS_NOTEBOOK = true

    // Response 
    c.String(http.StatusOK, "OK") 
}



