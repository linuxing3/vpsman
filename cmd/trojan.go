/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/linuxing3/vpsman/util"
	"github.com/spf13/cobra"
)

// trojanCmd represents the trojan command
var trojanCmd = &cobra.Command{
	Use:   "trojan",
	Short: "A brief description of your command",
	Long: ` About usage of using trojan. For example: 
Trojan is a CLI Command for Go that empowers proxy.
to quickly create a web tunnel.`,
	Run: func(cmd *cobra.Command, args []string) {
      trojanMenu()
	},
}

func trojanMenu() {
exit:
   for {
    fmt.Println()
    fmt.Print(util.Cyan("Please select command"))
    fmt.Println()
    loopMenu := []string{"Start", "Status", "Stop"}
    choice := util.LoopInput("Enter to Exit", loopMenu, false)
    switch choice {
        case 1:
            fmt.Println("trojan start")
            util.ExecCommand("systemctl start trojan")
        case 2:
            fmt.Println("trojan status")
            util.ExecCommand("systemctl status trojan")
        case 3:
            fmt.Println("trojan stop")
            util.ExecCommand("systemctl stop trojan")
        default:
            break exit
        }
   }
}
func init() {
	rootCmd.AddCommand(trojanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// trojanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// trojanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
