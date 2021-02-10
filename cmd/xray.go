package cmd

import (
	"fmt"

	"github.com/linuxing3/vpsman/util"
	"github.com/spf13/cobra"
)

// xrayCmd represents the xray command
var xrayCmd = &cobra.Command{
	Use:   "xray",
	Short: "Manage xray service on your vps",
	Long: ` About usage of using xray. For example: 
Xray is a CLI Command for Go that empowers proxy.
to quickly create a web tunnel.`,
	Run: func(cmd *cobra.Command, args []string) {
      xrayMenu()
	},
}

func xrayMenu() {
exit:
   for {
    fmt.Println()
    fmt.Print(util.Cyan("请选择"))
    fmt.Println()
    loopMenu := []string{"启动", "状态", "停止"}
    choice := util.LoopInput("回车退出", loopMenu, false)
    switch choice {
        case 1:
            fmt.Println("xray start")
            util.ExecCommand("systemctl start xray")
        case 2:
            fmt.Println("xray status")
            util.ExecCommand("systemctl status xray")
        case 3:
            fmt.Println("xray stop")
            util.ExecCommand("systemctl stop xray")
        default:
            break exit
        }
   }
}

func init() {
	rootCmd.AddCommand(xrayCmd)
}
