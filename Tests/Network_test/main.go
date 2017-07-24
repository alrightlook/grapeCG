// 2017-07-24
// 测试网络层基本收法行为
package main

import (
	"log"
	"runtime"

	logger "github.com/koangel/grapeNet/Logger"
	gnet "github.com/koangel/grapeNet/Net"

	"net/http"
	_ "net/http/pprof"
)

type userData struct {
	userName string
}

func CreateOwner() interface{} {
	return &userData{}
}

func HandleData(conn *gnet.TcpConn, ownerPak []byte) {
	logger.INFO("recv Data:%v - pak:%v", conn.SessionId, string(ownerPak))
}

func main() {
	// 设置并行运行
	runtime.GOMAXPROCS(runtime.NumCPU())

	curDir := logger.GetCurrentDirectory() + "/log"
	logger.BuildLogger(curDir, "normal.log")
	logger.INFO("Test Server Start...")
	cnet, err := gnet.NewTcpServer(":9923")
	if err != nil {
		logger.FLUSH()
		return
	}

	go func() {
		log.Println(http.ListenAndServe("localhost:7777", nil))
	}()

	cnet.Unpackage = gnet.DefaultLineData
	cnet.OnHandler = HandleData
	cnet.CreateUserData = CreateOwner

	cnet.Runnable()
}
