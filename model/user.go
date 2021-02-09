package model

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
	"github.com/linuxing3/vpsman/core"
	"github.com/linuxing3/vpsman/util"
)

// UserMenu 用户管理菜单
func UserMenu(dbPath string) {
	fmt.Println()
	menu := []string{"查询用户", "添加用户","更新用户","删除用户"}
	switch util.LoopInput("请选择: ", menu, false) {
	case 1:
		QueryAllUser(dbPath)
	case 2:
		AddUser(dbPath)
	case 3:
		UpdateUser(dbPath)
	case 4:
		DelUser(dbPath)
	}
}

// QueryAllUser 删除用户
func QueryAllUser(dbPath string) (*core.Sqlite, []*core.User ) {
	
	sqlite := core.NewSqlite(dbPath)
	userList, err := sqlite.QueryUsersORM()
	if err != nil {
		fmt.Print(err)
		return nil, nil
	}
	fmt.Println("Quering all users:")
	for i, k := range userList {
		fmt.Printf("%d.\n", i+1)
		fmt.Println("用户id: " + fmt.Sprint(k.ID))
		fmt.Println("用户名: " + k.Username)
	}
	return sqlite, userList
}


// AddUser 添加用户
func AddUser(dbPath string) error {

	randomUser := util.RandString(4)
	randomPass := util.RandString(8)
	inputUser := util.Input(fmt.Sprintf("生成随机用户名: %s, 使用直接回车, 否则输入自定义用户名: ", randomUser), randomUser)
	if inputUser == "admin" {
		fmt.Println(util.Yellow("不能新建用户名为'admin'的用户!"))
		return nil
	}
	inputPass := util.Input(fmt.Sprintf("生成随机密码: %s, 使用直接回车, 否则输入自定义密码: ", randomPass), randomPass) // originPass
	uuid := fmt.Sprintf("%s", uuid.New())
	fmt.Println(util.Yellow("[uuid]:" + uuid))

	sqlite := core.NewSqlite(dbPath)
	if sqlite.HasDuplicateUserORM(inputUser, inputPass) {
		fmt.Println("已存在这个用户名或密码的用户!")
		return nil
	}
	base64Pass := base64.StdEncoding.EncodeToString([]byte(inputPass)) // passwordShow
	// 创建Sqlite新用户
	if err := sqlite.CreateUserORM(uuid, inputUser, base64Pass, inputPass); err != nil {
		return err
	} 
	return nil

}

// DelUser 删除用户
func DelUser(dbPath string) error {
	
	sqlite, userList := QueryAllUser(dbPath)
	
	choice := util.LoopInput("请选择要删除的用户序号: ", userList, true)
	if choice == -1 {
		return nil
	}
	if err := sqlite.DeleteUserORM(fmt.Sprint(userList[choice-1].ID)); err != nil {
		return err
	}
	return nil
}

// UpdateUser 删除用户
func UpdateUser(dbPath string) error {
	
	sqlite, userList := QueryAllUser(dbPath)

	choice := util.LoopInput("请选择要修改的用户序号: ", userList, true)
	if choice == -1 {
		return nil
	}
	// Updating user information
	inputName := util.Input("请输入新名称:", "daniel")
	inputPass := util.Input("请输入密码", "000000")
	encryPass := sha256.Sum224([]byte(inputPass)) // %x => password
	base64Pass := base64.StdEncoding.EncodeToString([]byte(inputPass)) // passwordShow

	data := core.User{
		Username: inputName,
		Password: fmt.Sprintf("%x", encryPass),
		PasswordShow: base64Pass,
	}
	if err := sqlite.UpdateUserCondORM(fmt.Sprint(userList[choice-1].ID),&data); err != nil {
		return err
	}
	return nil
}