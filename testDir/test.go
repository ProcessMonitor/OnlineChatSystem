package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func LocalVarT() {
	x := 1
	fmt.Println(x) //1
	{
		fmt.Println(x) //1
		x := 2
		fmt.Println(x) //2
	}
	fmt.Println(x) //1

}
func MutiVarInitT() {
	const name, age = "Kim", 22
	s := fmt.Sprint(name, " is ", age, " years old.\n")
	io.WriteString(os.Stdout, s)
}

func StringArrT() {
	n := len("buf\n")
	msg := string("buf\n"[:n-1])
	print(msg)
	msg1 := string("buf\n")
	print(msg1)
}

func SombineStringArrayElements() {
	msgArray := []string{"/to", "用户名称", "消息内容", "格式发送"}
	privateMsgContentArray := msgArray[2:]
	privateMsgContent := strings.Join(privateMsgContentArray, " ")
	// 输出拼接好的字符串
	println(privateMsgContent)
}

func main() {
	//LocalVarT()
	//
	//MutiVarInitT()
	//
	//StringArrT()
	//ChannT()
	//mutiChannSelectT()
	//buffChanT()
	SombineStringArrayElements()
}
