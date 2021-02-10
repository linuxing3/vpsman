package controller

import (
	"fmt"
	"time"

	"github.com/linuxing3/vpsman/core"
	"github.com/linuxing3/vpsman/util"
)

// DefaultDbPath 数据库地址
var DefaultDbPath = "./vpsman.db"

// UserList 查询用户
func UserList(findUser string) *ResponseBody {

	var users []*core.User

	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)

	// 按用户名查询UserList
	sqlite := core.NewSqlite(DefaultDbPath)
	userWithName := core.User{
		Username: findUser,
	}
	if findUser != "" {
		if users, _ = sqlite.QueryUsersWithStructORM(&userWithName); len(users) > 0 {
			responseBody.Data = map[string]interface{}{
				"userList": users,
			}
		}
	} else {
		responseBody.Msg = "必须提供名称"
	}
	return &responseBody
}

// PageUserList 分页查询获取用户列表
func PageUserList(curPage int, pageSize int) *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	// connect db
	sqlite := core.NewSqlite(DefaultDbPath)
	// get pageData with curPage and pageSize
	pageData, err := sqlite.PageQueryUsersORM(curPage, pageSize)
	if err != nil {
		responseBody.Msg = err.Error()
		return &responseBody
	}
	responseBody.Data = map[string]interface{}{
		"pageData": pageData,
	}
	return &responseBody
}

// CreateUser 创建用户
func CreateUser(username string, password string) *ResponseBody {
	fmt.Printf("创建用户 %s, 密码是 %s", username, password)
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)

	if username == "admin" {
		responseBody.Msg = "不能创建用户名为admin的用户!"
		return &responseBody
	}
	// 创建普通用户
	sqlite := core.NewSqlite(DefaultDbPath)
	if sqlite.HasDuplicateUserORM(username, password) {
		responseBody.Msg = "已存在这个用户名或密码的用户!"
		return &responseBody
	}
	if err := sqlite.CreateUserORM("", username, password); err != nil {
		responseBody.Msg = err.Error()
	}
	return &responseBody
}

// UpdateUser 更新用户
func UpdateUser(id string, username string, password string) *ResponseBody {
	
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	
	if username == "admin" {
		responseBody.Msg = "不能更改用户名为admin的用户!"
		return &responseBody
	}
	
	sqlite := core.NewSqlite(DefaultDbPath)
	// 准备数据
	fmt.Printf("更新用户 %s, 新密码是 %s", username, password)
	encryPass, base64Pass := util.GenPass(password)
	data := core.User{
		Password: fmt.Sprintf("%x", encryPass),
		PasswordShow: base64Pass,
	}
	// 有id就用id更改
	if id == "" {
		if err := sqlite.UpdateUserCondORM(&core.User{Username: username}, &data); err != nil {
			responseBody.Msg = err.Error()
			return &responseBody
		}
		return &responseBody
	}
	// 没有id用姓名修改，不可改名
	if err := sqlite.UpdateUserByIdORM(id, &data); err != nil {
		responseBody.Msg = err.Error()
		return &responseBody
	}
	return &responseBody
}

// DelUser 删除用户
func DelUser(id string) *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	// 删除用户
	sqlite := core.NewSqlite(DefaultDbPath)
	err := sqlite.DeleteUserORM(id)
	if err != nil {
		responseBody.Msg = "Failed to delete"
	}
	return &responseBody
}

// SetExpire 设置用户过期
func SetExpire(id string, useDays uint) *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)

	sqlite := core.NewSqlite(DefaultDbPath)
	data := core.User{
		ExpiryDate: "",
	}
	if err := sqlite.UpdateUserByIdORM(id, &data); err != nil {
		responseBody.Msg = err.Error()
		return &responseBody
	}
	return &responseBody
}

// CancelExpire 取消设置用户过期
func CancelExpire(id string) *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
		
	sqlite := core.NewSqlite(DefaultDbPath)
	data := core.User{
		ExpiryDate: "1000000",
	}
	if err := sqlite.UpdateUserByIdORM(id, &data); err != nil {
		responseBody.Msg = err.Error()
		return &responseBody
	}
	return &responseBody
}
