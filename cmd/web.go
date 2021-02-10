package cmd

import (
	"crypto/sha256"
	"fmt"

	"github.com/linuxing3/vpsman/core"
	"github.com/linuxing3/vpsman/util"
	"github.com/linuxing3/vpsman/web"
	"github.com/spf13/cobra"
)

var (
	host string
	port int
	ssl  bool
)

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start Gin Web Server",
	Long: `A Web Server which offer UI for vpsman
and usage of using command. For example:

vpsman web -p 8080 --host 0.0.0.0 --ssl false
to quickly create a Web Gin application.`, 
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting web")
		startWebServer(host, port, ssl)
	},
}

// WebMenu get webmenu
func startWebServer(host string, port int, ssl bool) {
	web.Start(host, port, ssl)
}

// WebMenu web管理菜单
func webMenu() {
	fmt.Println()
	menu := []string{"重置web管理员密码", "启动web服务器"}
	switch util.LoopInput("请选择: ", menu, true) {
	case 1:
		ResetAdminPass()
	case 2:
		web.Start("0.0.0.0", 8888, false)
	}
}

// ResetAdminPass 重置管理员密码
func ResetAdminPass() {
	inputPass := util.Input("请输入admin用户密码: ", "")
	if inputPass == "" {
		fmt.Println("撤销更改!")
	} else {
		encryPass := sha256.Sum224([]byte(inputPass))
		err := core.SetValue("admin_pass", fmt.Sprintf("%x", encryPass))
		if err == nil {
			fmt.Println(util.Green("重置admin密码成功!"))
		} else {
			fmt.Println(err)
		}
	}
}

func init() {

	webCmd.Flags().StringVarP(&host, "host", "", "0.0.0.0", "web服务监听地址")
	webCmd.Flags().IntVarP(&port, "port", "p", 8888, "web服务启动端口")
	webCmd.Flags().BoolVarP(&ssl, "ssl", "", false, "web服务是否以https方式运行")

	rootCmd.AddCommand(webCmd)
}
