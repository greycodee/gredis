package cmd

import (
	"strconv"
)

var	crlf = []byte{0x0d,0x0a}

var cmdsMap = map[string]interface{}{
	"12":"2",
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
	for _,r := range cmd{
		if r == 32{
			// 空格
			continue
		}
		tmp = append(tmp, byte(r))
	}
	return result
}