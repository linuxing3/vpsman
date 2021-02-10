package cmd

import (
	"fmt"

	"github.com/linuxing3/vpsman/model"
	"github.com/linuxing3/vpsman/util"
	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manager users in sqlite database",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(dbPath)
		userMenu(dbPath)
	},
}

// UserMenu 用户管理菜单
func userMenu(dbPath string) {
	fmt.Println()
	menu := []string{"查询用户", "添加用户","更新用户","删除用户"}
	switch util.LoopInput("请选择: ", menu, false) {
	case 1:
		model.QueryAllUser(dbPath)
	case 2:
		model.AddUser(dbPath)
	case 3:
		model.UpdateUser(dbPath)
	case 4:
		model.DelUser(dbPath)
	}
}
func init() {
	userCmd.Flags().StringVarP(&dbPath, "db", "", "./vpsman.db", "数据库目录地址")
	rootCmd.AddCommand(userCmd)
}
