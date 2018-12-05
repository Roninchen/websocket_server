package impl

import (
	"github.com/gorilla/websocket"
	"sync"
	"github.com/pkg/errors"
)
var M = make(map[string]*websocket.Conn)
type Connection struct {
	wsConn *websocket.Conn
	inChan chan []byte
	outChan chan []byte
	closeChan chan byte

	sync.Mutex
	isClosed bool
}

func InitConnection(wsConn *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConn:wsConn,
		inChan:make(chan[]byte,1000),
		outChan:make(chan []byte,1000),
		closeChan:make(chan byte,1),
	}
	//读协程
	go conn.readLoop()
	//启动写携程
	go conn.writeLoop()
	return conn,nil
}
//API
func (conn *Connection)ReadMessage()(data []byte,err error)  {
	select {
		case data=<-conn.inChan:
		case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return 
}

func (conn *Connection)WriteMessage(data []byte)(err error)  {
	select {
		case conn.outChan<-data:
		case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

func (conn *Connection) Close() {
	//线程安全，可重入的close
	conn.wsConn.Close()
	//这一行代码只执行一次
	conn.Mutex.Lock()
	if !conn.isClosed {
		close(conn.closeChan)
		conn.isClosed=true
	}
	conn.Mutex.Unlock()
}
//内部实现
func (conn *Connection) readLoop() {
	var (
		data[]byte
		err error
	)
	for{
		if _,data,err =conn.wsConn.ReadMessage();err!=nil{
			goto ERR
		}
		select {
				case conn.inChan <- data:
				case <-conn.closeChan:
					//closeChan关闭的时候
				goto ERR
		}
	}
ERR:
	conn.Close()
}

func (conn *Connection)writeLoop(){
	var(
		data []byte
	)
	for{
		select {
			case data=<-conn.outChan:
			case <-conn.closeChan:
				goto ERR
		}
			if err := conn.wsConn.WriteMessage(websocket.TextMessage,data);err!=nil{
				goto ERR
			}



	}
ERR:
	conn.Close()
}