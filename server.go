package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	"ginserver/impl"
	"time"
	"fmt"
)
var (
	upgrader = websocket.Upgrader{
		//允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

)

func wsHandler(w http.ResponseWriter, r *http.Request){
	var(
		wsConn *websocket.Conn
		err error
		data []byte
		conn *impl.Connection
	)
	//Upgrade websocket
	if wsConn, err = upgrader.Upgrade(w, r, nil);err!=nil{
		return
	}

	if conn,err = impl.InitConnection(wsConn);err!=nil {
		goto ERR
	}

	impl.M[wsConn.RemoteAddr().String()]=wsConn
	for  k,_:=range impl.M {
		fmt.Println(k)
	}
	fmt.Println("=========================")
	go func() {
		var (
			err error
		)
		for {
			if err = conn.WriteMessage([]byte("heartbeat")); err != nil {
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()

	for{
		if data,err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		if err = conn.WriteMessage(data);err!=nil {
			goto ERR
		}
	}
ERR:
	defer conn.Close()
	
}
func main() {
	//http://localhost:7777//ws
	http.HandleFunc("/ws",wsHandler)
	http.ListenAndServe("0.0.0.0:7777",nil)
}