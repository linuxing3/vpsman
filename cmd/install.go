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
	"net"

	"github.com/gobuffalo/packr/v2"
	"github.com/linuxing3/vpsman/core"
	"github.com/linuxing3/vpsman/util"
	"github.com/spf13/cobra"
)

type Mysql struct {
	Enabled    bool   `json:"enabled"`
	ServerAddr string `json:"server_addr"`
	ServerPort int    `json:"server_port"`
	Database   string `json:"database"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Cafile     string `json:"cafile"`
}

var (
	dockerInstallUrl1 = "https://get.docker.com"
	dockerInstallUrl2 = "https://git.io/docker-install"
	XrayDbDockerRun   = "docker run --name xray-mariadb --restart=always -p %d:3306 -v /home/mariadb/xray:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=%s -e MYSQL_ROOT_HOST=%% -e MYSQL_DATABASE=%s -d mariadb:10.2"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Common tools or modules",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("install called")
		installMenu()
	},
}

func installMenu() {
	fmt.Println()
	menu := []string{"更新xray", "证书申请", "安装mysql", "安装sqlite"}
	switch util.LoopInput("请选择: ", menu, true) {
	case 1:
		InstallXray()
	case 2:
		InstallTls()
	case 3:
		InstallMysql(XrayDbDockerRun, "xray")
	case 4:
		InstallSqliteBin()
	default:
		return
	}
}

// InstallDocker 安装docker
func InstallDocker() {
	if !util.CheckCommandExists("docker") {
		util.RunWebShell(dockerInstallUrl1)
		if !util.CheckCommandExists("docker") {
			util.RunWebShell(dockerInstallUrl2)
		} else {
			util.ExecCommand("systemctl enable docker")
			util.ExecCommand("systemctl start docker")
		}
		fmt.Println()
	}
}

// InstallSqliteBin 安装sqlite3的可执行文件
func InstallSqliteBin() {
	util.ExecCommand("apt install -y sqlite3")
}

// InstallXray 安装xray
func InstallXray() {
	fmt.Println()
	box := packr.New("xray-install", "../asset")
	data, err := box.FindString("xray-install.sh")
	if err != nil {
		fmt.Println(err)
	}
	util.ExecCommand(data)
	util.OpenPort(443)
	util.ExecCommand("systemctl restart xray")
	util.ExecCommand("systemctl enable xray")
}

// InstallTls 安装证书
func InstallTls() {
	domain := ""
	fmt.Println()
	choice := util.LoopInput("请选择使用证书方式: ", []string{"Let's Encrypt 证书", "自定义证书路径"}, true)
	if choice < 0 {
		return
	} else if choice == 1 {
		localIP := util.GetLocalIP()
		fmt.Printf("本机ip: %s\n", localIP)
		for {
			domain = util.Input("请输入申请证书的域名: ", "")
			ipList, err := net.LookupIP(domain)
			fmt.Printf("%s 解析到的ip: %v\n", domain, ipList)
			if err != nil {
				fmt.Println(err)
				fmt.Println("域名有误,请重新输入")
				continue
			}
			checkIp := false
			for _, ip := range ipList {
				if localIP == ip.String() {
					checkIp = true
				}
			}
			if checkIp {
				break
			} else {
				fmt.Println("输入的域名和本机ip不一致, 请重新输入!")
			}
		}
		util.InstallPack("socat")
		if !util.IsExists("/root/.acme.sh/acme.sh") {
			util.RunWebShell("https://get.acme.sh")
		}
		util.ExecCommand("systemctl stop xray-web")
		util.OpenPort(80)
		util.ExecCommand(fmt.Sprintf("bash /root/.acme.sh/acme.sh --issue -d %s --debug --standalone --keylength ec-256", domain))
		crtFile := "/root/.acme.sh/" + domain + "/fullchain.cer"
		keyFile := "/root/.acme.sh/" + domain + "/" + domain + ".key"
		// 写入证书到xray配置文件
		core.WriteTls(crtFile, keyFile, domain)
	} else if choice == 2 {
		crtFile := util.Input("请输入证书的cert文件路径: ", "")
		keyFile := util.Input("请输入证书的key文件路径: ", "")
		if !util.IsExists(crtFile) || !util.IsExists(keyFile) {
			fmt.Println("输入的cert或者key文件不存在!")
		} else {
			domain = util.Input("请输入此证书对应的域名: ", "")
			if domain == "" {
				fmt.Println("输入域名为空!")
				return
			}
			core.WriteTls(crtFile, keyFile, domain)
		}
	}
	util.ExecCommand("systemctl restart trojan-web")
	fmt.Println()
}

// InstallMysql InstallMysql with docker
func InstallMysql(dockerCommand string, database string) {
	var (
		server   string
		username string
		mysql    Mysql
	)
	server = "127.0.0.1"
	username = "root"
	fmt.Println()
	mysql = Mysql{ServerAddr: server, ServerPort: util.RandomPort(), Password: util.RandString(5), Username: username, Database: database}
	// install docker
	InstallDocker()
	// 显示说明：链接并创建一个xray的数据库
	fmt.Println(fmt.Sprintf(dockerCommand, mysql.ServerPort, mysql.Password, database))
	if util.CheckCommandExists("setenforce") {
		util.ExecCommand("setenforce 0")
	}
	util.OpenPort(mysql.ServerPort)
	// 执行命令: 创建trojan数据库
	util.ExecCommand(fmt.Sprintf(dockerCommand, mysql.ServerPort, mysql.Password, database))

	fmt.Println("mariadb启动成功!")
}

func init() {
	rootCmd.AddCommand(installCmd)
}
