package cmd

import (
	"strconv"
	"strings"
)

var	crlf = []byte{0x0d,0x0a}

var cmdMap = map[string]func(c string) []byte {
	"set":set(),
	"get":get(),
}

func GetCmdByte(c string) ([]byte,bool,string) {
	//s := strings.Split(c," ")
	//if f,ok := cmdMap[s[0]]; ok{
	//	return f(c),nil
	//}else{
	//	return nil,errors.New("cmd error")
	//}
	s := strings.Split(c," ")
	return buildCommand(readCmdString(c)...),s[0]=="select",s[len(s)-1]
}

func  buildCommand(cmd ...string) []byte {
	cmdLen := strconv.Itoa(len(cmd))
	cmdBytes := make([]byte,0)
	for _,v := range cmd{
		vv := []byte(v)
		cmdBytes = append(cmdBytes, 0x24)
		cmdBytes = append(cmdBytes, []byte(strconv.Itoa(len(vv)))...)
		cmdBytes = append(cmdBytes, crlf...)
		cmdBytes = append(cmdBytes, vv...)
		cmdBytes = append(cmdBytes, crlf...)
	}
	request := commonRequest(cmdLen,cmdBytes)
	return request
}

func  commonRequest(cmdLen string,cmd []byte) []byte {
	request := make([]byte,0)
	request = append(request, 0x2a)
	request = append(request, []byte(cmdLen)...)
	request = append(request, crlf...)
	request = append(request, cmd...)
	return request
}

func readCmdString(cmd string) []string {
	result := make([]string,0)
	tmp := make([]byte,0)

	quot := false
	for _,r := range []byte(cmd){
		// 34 双引号
		if r == 34 {
			quot = !quot
			if !quot{
				if len(tmp)>0 {
					result = append(result, string(tmp))
					tmp = make([]byte,0)
				}
				continue
			}
			continue
		}
		if !quot {
			if r == 32{
				// 空格
				if len(tmp)>0 {
					result = append(result, string(tmp))
					tmp = make([]byte,0)
				}
				continue
			}
		}
		tmp = append(tmp, r)
	}
	if len(tmp)>0 {
		result = append(result, string(tmp))
		tmp = make([]byte,0)
	}
	return result
}