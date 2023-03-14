package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	//用户管理的channel 不停检查
	User_Message_Chann chan string
	//客户端的触手 句柄
	connection net.Conn

	//当前用户属于哪个Server （双向绑定
	server *Server
}

// user api
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		// 用户端的消息管道
		User_Message_Chann: make(chan string),
		// 该User与服务器的链接 （客户端
		connection: conn,
		// 该User隶属的服务器
		server: server,
	}
	//创建 goroutine去监听 channel
	go user.ListenChannel()

	return user
}

// 用户发消息给服务器的方法
func (this *User) ListenChannel() {
	for {
		//一直遍历 检查有无新消息
		msg := <-this.User_Message_Chann
		// 把消息 写成二进制 发给客户端and发往服务器
		this.connection.Write([]byte(msg + "\n"))
	}
}

// 上线方法
func (this *User) Online() {
	//add to uer online map
	this.server.MapLock.Lock()
	this.server.OnlineUserMap[this.Name] = this
	this.server.MapLock.Unlock()
	// broadcast online message to all users
	this.server.BroadCast(this, "我上线啦")
}

// 下线方法
func (this *User) Offline() {
	//delete from uer online map
	this.server.MapLock.Lock()
	delete(this.server.OnlineUserMap, this.Name)
	this.server.MapLock.Unlock()
	// broadcast online message to all users
	this.server.BroadCast(this, "我下了886")
}

func (this *User) sendToClient(msg string) {
	this.connection.Write([]byte(msg))
}

// 给当前用户的客户端发消息
func (this *User) DoMessage(msg string) {
	if msg == "/who" {
		//查看谁在线
		this.server.MapLock.Lock()
		for _, user := range this.server.OnlineUserMap {
			onlineUserNotify := "[" + user.Name + "]is online!\n"
			//用户给客户端发消息
			this.sendToClient(onlineUserNotify)
		}
		this.server.MapLock.Unlock()
	} else if len(msg) > 6 && msg[:6] == "/cname" {
		// pattern : /cname xxx
		hasSpace := strings.Contains(msg, " ")
		if hasSpace == true {
			//切割string
			newName := strings.Split(msg, " ")[1]
			_, ok := this.server.OnlineUserMap[newName]
			if ok {
				this.sendToClient("已被他人使用!\n")
			} else {
				this.server.MapLock.Lock()
				delete(this.server.OnlineUserMap, this.Name)
				this.server.OnlineUserMap[newName] = this
				this.server.MapLock.Unlock()
				this.Name = newName
				this.sendToClient("名称已更改!\n")
			}
		}
	} else if len(msg) > 4 && msg[:3] == "/to" {
		//pattern : /to remoteUser content
		hasSpace := strings.Contains(msg, " ")
		if hasSpace == true {
			// 切割string
			msgArray := strings.Split(msg, " ")
			toName := msgArray[1]
			if toName == "" {
				this.sendToClient("消息格式不正确 /to 用户名称 消息内容 格式发送")
			}
			// 从服务器找到要发送信息的User
			remoteUser, ok := this.server.OnlineUserMap[toName]
			// 如果不存在
			if !ok {
				this.sendToClient(toName + "此用户不存在\n")
				return
			}
			// 得到信息本体数组
			privateMsgContentArray := msgArray[2:]
			// 重新分割后面的空格
			privateMsgContent := strings.Join(privateMsgContentArray, " ")
			// 如果消息为空
			if privateMsgContent == "" {
				this.sendToClient("消息不能为空\n")
				return
			}
			// send
			remoteUser.sendToClient(this.Name + " 对你说: " + privateMsgContent)
		}
	} else {
		//normal broadCast method
		this.server.BroadCast(this, msg)
	}
}
