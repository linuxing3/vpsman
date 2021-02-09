package controller

import (
	"fmt"
	"time"

	"github.com/linuxing3/vpsman/util"
	websocket "github.com/linuxing3/vpsman/util"

	"github.com/gin-gonic/gin"
)

// Start 启动xray
func Start() *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	util.ExecCommand("systemctl start xray")
	return &responseBody
}

// Stop 停止xray
func Stop() *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	util.ExecCommand("systemctl stop xray")
	return &responseBody
}

// Restart 重启xray
func Restart() *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	util.ExecCommand("systemctl restart xray")
	return &responseBody
}

// Log 通过ws查看xray实时日志
func Log(c *gin.Context) {
	var (
		wsConn *websocket.WsConnection
		err    error
	)
	if wsConn, err = websocket.InitWebsocket(c.Writer, c.Request); err != nil {
		fmt.Println(err)
		return
	}
	defer wsConn.WsClose()
	param := c.DefaultQuery("line", "300")
	if param == "-1" {
		param = "--no-tail"
	} else {
		param = "-n " + param
	}
}
