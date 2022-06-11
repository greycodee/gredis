package core

import (
	"bufio"
	"encoding/json"
	"net"
	"strconv"
	"strings"
)

type RedisClient struct {
	conn net.Conn
	pipeline []byte

	br  *bufio.Reader
	bw 	*bufio.Writer
}

func (c *RedisClient) Open(addr string)  (err error) {
	c.conn, err = net.Dial("tcp",addr)
	c.br = bufio.NewReader(c.conn)
	c.bw = bufio.NewWriter(c.conn)
	return
}

func (c RedisClient) handleResp()  (resp []byte) {
	//flag := 0x11
	flag, err := c.br.Peek(1)

	if err != nil {
		return
	}
	switch flag[0] {
	case 0x2b:
		// + Simple Strings
		s := c.parseSimpleString()
		resp = []byte(s)
		break
	case 0x2d:
		// - Errors
		e := c.parseErrors()
		resp = []byte(e)
		break
	case 0x3a:
		// : Integers
		i := c.parseIntegers()
		//resp = make([]byte, 8)
		//binary.LittleEndian.PutUint64(resp,uint64(i))
		resp = []byte(strconv.Itoa(int(i)))
		break
	case 0x24:
		// $ Bulk Strings
		s := c.parseBulkStrings()
		resp = []byte(s)
		break
	case 0x2a:
		// * Arrays
		resp = []byte(c.parseArrays())
		break
	}
	return
}

func (c *RedisClient) ExecCMD(cmd ...string)  (resp []byte){
	_, err := c.bw.Write(buildCommand(cmd...))
	if err != nil {
		return
	}
	err = c.bw.Flush()
	if err != nil {
		return
	}
	return c.readConn()
}

func (c *RedisClient) ExecCMDByte(cmdByte []byte)  (resp []byte){
	_, err := c.bw.Write(cmdByte)
	if err != nil {
		return
	}
	err = c.bw.Flush()
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

func (c *RedisClient) readConn() ([]byte) {
	c.br.Reset(c.conn)
	return c.handleResp()
}

func (c *RedisClient) readPipelineConn() []byte {
	resp := make([]byte,1024)
	n, _ := c.conn.Read(resp)
	resp = resp[:n]
	index := 0
	result := make([]string,0)
	for index < n{
		result = append(result, string(c.handleResp()))
	}
	marshal, err := json.Marshal(result)
	if err != nil {
		return nil
	}
	return marshal
}

func (c RedisClient) parseSimpleString() string {
	return c.simpleParse()
}
func (c RedisClient) parseIntegers() int64 {
	parseInt, _ := strconv.ParseInt(string(c.readLine()), 10, 64)
	return parseInt
}
func (c RedisClient) parseErrors() string{
	return string(c.readLine())
}

func (c RedisClient) simpleParse() string  {
	return string(c.readLine())
}

func (c RedisClient) readLine() []byte {
	line, _, err := c.br.ReadLine()
	if err != nil {
		return nil
	}
	return line[1:]
}

func (c RedisClient) parseBulkStrings() string {
	// 获取字节长度
	strLen, _, err := c.br.ReadLine()
	if err != nil {
		return ""
	}
	l, _ := strconv.ParseInt(string(strLen[1:]),10,64)
	if l == -1 {
		return ""
	}
	str := make([]byte,l)
	read, err := c.br.Read(str)
	c.br.ReadLine()
	if err != nil {
		return ""
	}
	return string(str[:read])
}
func (c RedisClient) parseArrays() string  {
	strLen, _, err := c.br.ReadLine()
	if err != nil {
		return ""
	}
	l, _ := strconv.ParseInt(string(strLen[1:]),10,64)
	//result := make([]string,l)
	arrays := strings.Builder{}
	for i:=int64(0);i<l;i++ {
		arrays.Write(c.handleResp())
		arrays.Write(crlf)
		//result = append(result, string(r))
	}
	return arrays.String()
}

func (c *RedisClient) pipelineCMDAdd(cmd ...string)  {
	c.pipeline= append(c.pipeline, buildCommand(cmd...)...)
}
func (c *RedisClient) pipelineExec() []byte {
	_, err := c.conn.Write(c.pipeline)
	if err != nil {
		return nil
	}
	return c.readPipelineConn()
}
