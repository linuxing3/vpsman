package model

import (
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
	"github.com/linuxing3/vpsman/core"
	"github.com/linuxing3/vpsman/util"
)

// UserMenu 用户管理菜单
func UserMenu(dbPath string) {
	fmt.Println()
	menu := []string{"新增用户", "删除用户"}
	switch util.LoopInput("请选择: ", menu, false) {
	case 1:
		AddUser(dbPath)
	case 2:
		DelUser(dbPath)
	}
}

// AddUser 添加用户
func AddUser(dbPath string) {

	randomUser := util.RandString(4)
	randomPass := util.RandString(8)
	inputUser := util.Input(fmt.Sprintf("生成随机用户名: %s, 使用直接回车, 否则输入自定义用户名: ", randomUser), randomUser)
	if inputUser == "admin" {
		fmt.Println(util.Yellow("不能新建用户名为'admin'的用户!"))
		return
	}
	// 1. uuid，用于xray
	uuid := fmt.Sprintf("%s", uuid.New())
	fmt.Println(util.Yellow("[uuid]:" + uuid))

	// 2. 生成随机密码，通过密码获取用户，存在报错
	inputPass := util.Input(fmt.Sprintf("生成随机密码: %s, 使用直接回车, 否则输入自定义密码: ", randomPass), randomPass)
	base64Pass := base64.StdEncoding.EncodeToString([]byte(inputPass))

	// 创建Sqlite新用户
	sqlite := core.NewSqlite(dbPath)
	if err := sqlite.CreateUserORM(uuid, inputUser, base64Pass, inputPass); err != nil {
		fmt.Println("新增Sqlite用户成功!")
		fmt.Println("")
	} else {
		fmt.Println(err)
	}

}

// DelUser 删除用户
func DelUser(dbPath string) {
	
	sqlite := core.NewSqlite(dbPath)
	userList, err := sqlite.QueryUsersORM()
	if err != nil {
		fmt.Print(err)
	}
	choice := util.LoopInput("请选择要删除的用户序号: ", userList, true)
	if choice == -1 {
		return
	}
	if err := sqlite.DeleteUserORM(fmt.Sprint(userList[choice-1].ID)); err != nil {
		fmt.Println("删除Sqlite用户成功!")
		fmt.Println("")
	} else {
		fmt.Println(err)
	}
}

