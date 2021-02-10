package cmd

import (
	"fmt"

	"github.com/linuxing3/vpsman/util"
	"github.com/spf13/cobra"
)

// caddyCmd represents the caddy command
var caddyCmd = &cobra.Command{
	Use:   "caddy",
	Short: "A brief description of your command",
	Long: `About usage of using caddy. For example: 
	Caddy is a CLI Command for Go that empowers web.
	server with https support.`,
	Run: func(cmd *cobra.Command, args []string) {
		caddyMenu()
	},
}

func caddyMenu() {
	exit:
		 for {
			fmt.Println()
			fmt.Print(util.Cyan("请选择"))
			fmt.Println()
			loopMenu := []string{"启动", "状态", "停止"}
			choice := util.LoopInput("回车退出", loopMenu, false)
			switch choice {
					case 1:
							fmt.Println("caddy start")
							util.ExecCommand("systemctl start caddy")
					case 2:
							fmt.Println("caddy status")
							util.ExecCommand("systemctl status caddy")
					case 3:
							fmt.Println("caddy stop")
							util.ExecCommand("systemctl stop caddy")
					default:
							break exit
					}
		 }
	}

func init() {
	rootCmd.AddCommand(caddyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// caddyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// caddyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
