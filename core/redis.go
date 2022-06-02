package core

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

type RedisClient struct {
	conn net.Conn
	pipeline []byte
}

func (c *RedisClient) Open(addr string)  (err error) {
	c.conn, err = net.Dial("tcp",addr)
	return
}

var	crlf = []byte{0x0d,0x0a}

func (c RedisClient) handleResp(data []byte)  (resp []byte,endPoint int64){
	//flag := 0x11
	switch data[0] {
	case 0x2b:
		// + Simple Strings
		s,end := c.parseSimpleString(data)
		resp = []byte(s)
		endPoint = end
		break
	case 0x2d:
		// - Errors
		e,end := c.parseErrors(data)
		resp = []byte(e)
		endPoint = end
		break
	case 0x3a:
		// : Integers
		i,end := c.parseIntegers(data)
		//resp = make([]byte, 8)
		//binary.LittleEndian.PutUint64(resp,uint64(i))
		resp = []byte(strconv.Itoa(int(i)))
		endPoint = end
		break
	case 0x24:
		// $ Bulk Strings
		s,end := c.parseBulkStrings(data)
		resp = []byte(s)
		endPoint = end
		break
	case 0x2a:
		// * Arrays
		arrays,end := c.parseArrays(data)
		resp,_ = json.Marshal(arrays)
		endPoint = end
		break
	}
	return
}

func (c *RedisClient) ExecCMD(cmd ...string)  (resp []byte){
	_, err := c.conn.Write(c.buildCommand(cmd...))
	if err != nil {
		return
	}
	return c.readConn()
}

func (c *RedisClient) ExecCMDByte(cmdByte []byte)  (resp []byte){
	_, err := c.conn.Write(cmdByte)
	if err != nil {
		return
	}
	return c.readConn()
}

func (c *RedisClient) Shutdown()  {
	err := c.conn.Close()
	if err != nil {
		panic("redis connect close failed!")
		return
	}
}

func (c RedisClient) buildCommand(cmd ...string) []byte {
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
	request := c.commonRequest(cmdLen,cmdBytes)
	return request
}

func (c *RedisClient) commonRequest(cmdLen string,cmd []byte) []byte {
	request := make([]byte,0)
	request = append(request, 0x2a)
	request = append(request, []byte(cmdLen)...)
	request = append(request, crlf...)
	request = append(request, cmd...)
	return request
}


func (c *RedisClient) readConn() []byte {
	resp := make([]byte,1024)
	n, _ := c.conn.Read(resp)
	r,_ := c.handleResp(resp[:n])
	return r
}

func (c *RedisClient) readPipelineConn() []byte {
	resp := make([]byte,1024)
	n, _ := c.conn.Read(resp)
	resp = resp[:n]
	index := 0
	result := make([]string,0)
	for index < n{
		r,end := c.handleResp(resp[index:])
		result = append(result, string(r))
		index += int(end)+1
	}
	marshal, err := json.Marshal(result)
	if err != nil {
		return nil
	}
	return marshal
}

func (c RedisClient) parseSimpleString(data []byte) (s string,endIndex int64) {
	return c.simpleParse(data)
}
func (c RedisClient) parseIntegers(data []byte) (i int64,endIndex int64) {
	d,endPoint := c.simpleParse(data)
	parseInt, _ := strconv.ParseInt(d, 10, 64)
	return parseInt,endPoint
}
func (c RedisClient) parseErrors(data []byte) (s string,endIndex int64){
	return c.simpleParse(data)
}

func (c RedisClient) simpleParse(data []byte) (string,int64)  {
	s := make([]byte,0)
	for i,v :=range data[1:]{
		if v == 0x0d {
			return string(s), int64(i + 2)
		}
		s = append(s, v)
	}
	return "",0
}

func (c RedisClient) parseBulkStrings(data []byte) (s string,endIndex int64) {
	lenData := make([]byte,0)
	data = data[1:]
	index := 0
	for  i,v:= range data{
		if v ==0x0d {
			index = i+2
			break
		}
		lenData = append(lenData, v)
	}
	l, _ := strconv.ParseInt(string(lenData),10,64)
	if l == -1 {
		return *new(string), int64(index)
	}else if l == 0 {
		return "", int64(index + 2)
	}else {
		endIndex = int64(index)+l
		s = string(data[index:endIndex])
		endIndex+=2
		return
	}
}
func (c RedisClient) parseArrays(data []byte)([]string,int64)  {
	data = data[1:]
	lenData := make([]byte,0)
	index := 0
	for  i,v:= range data{
		if v ==0x0d {
			index = i+2
			break
		}
		lenData = append(lenData, v)
	}
	l, _ := strconv.ParseInt(string(lenData),10,64)
	fmt.Println(l)
	result := make([]string,0)
	for index<len(data) {
		resp,end := c.handleResp(data[index:])
		result = append(result, string(resp))
		index += int(end)+1
	}
	return result, int64(index)
}

func (c *RedisClient) pipelineCMDAdd(cmd ...string)  {
	c.pipeline= append(c.pipeline, c.buildCommand(cmd...)...)
}
func (c *RedisClient) pipelineExec() []byte {
	_, err := c.conn.Write(c.pipeline)
	if err != nil {
		return nil
	}
	return c.readPipelineConn()
}
