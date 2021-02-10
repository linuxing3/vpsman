package cmd

import (
	"fmt"

	"github.com/linuxing3/vpsman/util"
	"github.com/spf13/cobra"
)

// nginxCmd represents the nginx command
var nginxCmd = &cobra.Command{
	Use:   "nginx",
	Short: "A brief description of your command",
	Long: ` About usage of using nginx. For example: 
Nginx is a CLI Command for Go that empowers proxy.
to quickly create a web tunnel.`,
	Run: func(cmd *cobra.Command, args []string) {
      fmt.Println(args)
      nginxMenu()
	},
}

// NginxMenu 控制NginxMenu
func nginxMenu() {
exit:
   for {
    fmt.Println()
    fmt.Print(util.Cyan("Please select command"))
    fmt.Println()
    loopMenu := []string{"Start", "Status", "Stop"}
    choice := util.LoopInput("Enter to Exit", loopMenu, false)
    switch choice {
        case 1:
            fmt.Println("nginx start")
            util.ExecCommand("systemctl start nginx")
        case 2:
            fmt.Println("nginx status")
            util.ExecCommand("systemctl status nginx")
        case 3:
            fmt.Println("nginx stop")
            util.ExecCommand("systemctl stop nginx")
        default:
            break exit
        }
   }
}
func init() {
	rootCmd.AddCommand(nginxCmd)
	nginxCmd.Flags().BoolP("toggle", "t", false, "Toggle nginx service")
}
