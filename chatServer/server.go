package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

/**
本文件完成server 端的基本构建
*/
type Server struct {
	Ip   string
	Port int
	// online user list
	OnlineUserMap map[string]*User
	//read and write mutex
	MapLock sync.RWMutex
	// broadcast channel
	Server_Message_Chann chan string
}

// create server interface
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:                   ip,
		Port:                 port,
		OnlineUserMap:        make(map[string]*User),
		Server_Message_Chann: make(chan string),
	}
	return server
}

//User上下线的广播方法
func (this *Server) BroadCast(user *User, msg string) {
	sendmsg := "[" + user.Addr + "]" +
		user.Name + "say:" + msg
	this.Server_Message_Chann <- sendmsg
}

//配合goroutine 处理当前业务
func (this *Server) Handler(conn net.Conn) {

	fmt.Println("user connected success")
	user := NewUser(conn, this)
	user.Online()
	//online status
	isLive := make(chan bool)
	////add to uer online map
	//this.MapLock.Lock()
	//this.OnlineUserMap[user.Name] = user
	//this.MapLock.Unlock()
	//// broadcast online message to all users
	//this.BroadCast(user, "我上线啦")

	//receive message from user check if there has broadcast message
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn read err : ", err)
				return
			}
			//get User info (User 类加了 \n 这里要去除)
			msg := string(buf[:n-1])

			user.DoMessage(msg)
			// keep live
			isLive <- true
		}
	}()
	//block this handler
	for {
		select {
		case <-isLive:
			// 当前用户活跃重置定时器
			// 不做任何事情 激活select 更新定时器
		case <-time.After(time.Second * 30):
			//超时处理
			user.sendToClient("超时下线")
			//强制关闭

			close(user.User_Message_Chann)
			conn.Close()

			return // runtime.goExit()
		}
	}
}

//启动服务器接口
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp4",
		fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.listen err :", err)
		return
	}
	//close listen socket
	defer listener.Close()

	// 启动监听server message的goroutine
	go this.ListenServerMessage()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		//do handler //处理当前业务
		go this.Handler(conn)

	}
	//close listen socket
	listener.Close()
}

//监听广播 goroutine 有消息就发送给全部的在线User
func (this *Server) ListenServerMessage() {
	for {
		msg := <-this.Server_Message_Chann
		//message send
		this.MapLock.Lock()
		//  key ,value
		for _, cli := range this.OnlineUserMap {
			// 服务端广播
			cli.User_Message_Chann <- msg
		}
		this.MapLock.Unlock()
	}
}
