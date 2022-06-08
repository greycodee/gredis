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
	//p1 := tea.NewProgram(views.InitialModel(),tea.WithAltScreen())
	//if err := p1.Start(); err != nil {
	//	fmt.Printf("could not start program: %s\n", err)
	//	os.Exit(1)
	//}

	//block := lipgloss.PlaceHorizontal(10, lipgloss.Center, "asd")
	//
	//fmt.Printf("%s\n", block)

	// Tview
	//app := tview.NewApplication()
	//err := app.SetRoot(views.DrawFlex(),true).EnableMouse(true).Run()
	//if err != nil {
	//	return
	//}
	//fmt.Println(views.KeyType(6))
	//in := &views.KeyInfo{
	//	1,
	//	"123",
	//	views.STRING,
	//}

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
		loginResp,_ := redisClient.ExecCMD("auth",*passwd)
		if string(loginResp)!="OK" {
			// 登陆失败
			panic("password error!")
		}
	}

	// 选择 db
	resp,_ := redisClient.ExecCMD("select",*db)
	if string(resp)!="OK" {
		// 选择db失败
		panic("select database failed!")
	}
	return *redisClient
}
