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
