package controller

import (
	"encoding/base64"
	"time"

	"github.com/linuxing3/vpsman/core"
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
	cond := struct{
		Username string
	}{
		Username: findUser,
	}
	if findUser != "" {
		if users, _ = sqlite.QueryUsersWhereORM(cond); len(users) > 0 {
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
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)

	if username == "admin" {
		responseBody.Msg = "不能创建用户名为admin的用户!"
		return &responseBody
	}
	// 创建普通用户
	sqlite := core.NewSqlite(DefaultDbPath)
	// 1.解码密码作为origPass
	pass, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		responseBody.Msg = "Base64解码失败: " + err.Error()
		return &responseBody
	}
	// 2. 用加密过的密码查询是否重复
	cond := struct{
		Password string
	}{
		Password: password,
	}
	if users, _ := sqlite.QueryUsersWhereORM(cond); len(users) > 0 {
		responseBody.Msg = "已存在密码为: " + string(pass) + " 的用户!"
		return &responseBody
	}
	// 3. 创建普通用户
	if err := sqlite.CreateUserORM("", username, password, string(pass)); err != nil {
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
	// 使用id查询用户
	sqlite := core.NewSqlite(DefaultDbPath)
	foundUser, err := sqlite.QueryUserORM(id)
	if err != nil {
		responseBody.Msg = err.Error()
		return &responseBody
	}
	// 1. 检查姓名是否重复
	if foundUser.Username != username {
		cond := struct{
			UserName string
		}{
			UserName: username,
		}
		if users, _ := sqlite.QueryUsersWhereORM(cond); len(users) != 0 {
			responseBody.Msg = "已存在用户名为: " + username + " 的用户!"
			return &responseBody
		}
	}
	// 2. 检查密码是否重复
	if foundUser.Password != password {
		cond := struct{
			Password string
		}{
			Password: password,
		}
		if users, _ := sqlite.QueryUsersWhereORM(cond); len(users) != 0 {
			responseBody.Msg = "已存在密码的用户!"
			return &responseBody
		}
	}
	// 3. 解码密码作为origPass
	pass, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		responseBody.Msg = "Base64解码失败: " + err.Error()
		return &responseBody
	}
	// 4. 更新用户信息
	if err := sqlite.UpdateUserORM(id, username, password, string(pass)); err != nil {
		responseBody.Msg = err.Error()
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
	return &responseBody
}

// CancelExpire 取消设置用户过期
func CancelExpire(id string) *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	return &responseBody
}
