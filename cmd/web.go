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
		webMenu(host, port, ssl)
	},
}

// WebMenu get webmenu
func webMenu(host string, port int, ssl bool) {
	web.Start(host, port, ssl)
}


func init() {

	webCmd.Flags().StringVarP(&host, "host", "", "0.0.0.0", "web服务监听地址")
	webCmd.Flags().IntVarP(&port, "port", "p", 8888, "web服务启动端口")
	webCmd.Flags().BoolVarP(&ssl, "ssl", "", false, "web服务是否以https方式运行")

	rootCmd.AddCommand(webCmd)
}
