package main

import (
	"flag"
	"github.com/greycodee/gredis/core"
	"github.com/greycodee/gredis/views"
)

var (
	addr = flag.String("h","127.0.0.1","Redis Server Hosts")
	port = flag.String("p","6379","Redis Server Port")
	db = flag.String("d","0","Redis Server Databases")
	passwd = flag.String("auth","","Redis Server Password")
)

func init()  {
	flag.Parse()
}

func main() {
	// 连接 Redis 服务器
	conn := login()
	redisServer := &views.RedisServer{
		Addr: *addr,
		Port: *port,
		DB:   *db,
		Conn:  conn,
	}

	tui := views.TUI{RedisServer: *redisServer}
	tui.StartTUI()
}

func login() core.RedisClient {
	redisClient := &core.RedisClient{}
	err := redisClient.Open(*addr + ":" + *port)
	if err != nil {
		panic(err)
	}
	// 登陆认证 auth
	if *passwd != ""{
		loginResp := redisClient.ExecCMD("auth",*passwd)
		if string(loginResp)!="OK" {
			// 登陆失败
			panic("password error!")
		}
	}

	// 选择 db
	resp := redisClient.ExecCMD("select",*db)
	if string(resp)!="OK" {
		// 选择db失败
		panic("select database failed!")
	}
	return *redisClient
}
