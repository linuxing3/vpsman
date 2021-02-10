package cmd

import (
	"fmt"

	"github.com/linuxing3/vpsman/model"
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

// UserMenu get model for UserMenu
func userMenu(dbPath string) {
	model.UserMenu(dbPath)
}

func init() {
	webCmd.Flags().StringVarP(&dbPath, "db", "", "./vpsman.db", "数据库目录地址")
	rootCmd.AddCommand(userCmd)
}
