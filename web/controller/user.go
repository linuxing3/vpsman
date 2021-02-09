package controller

import (
	"encoding/base64"
	"time"

	"github.com/linuxing3/vpsman/core"
)

// DefaultDbPath 数据库地址
var DefaultDbPath = "./vpsman.db"

// UserList 获取用户列表
func UserList(findUser string) *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	responseBody.Data = map[string]interface{}{
	}
	return &responseBody
}

// PageUserList 分页查询获取用户列表
func PageUserList(curPage int, pageSize int) *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	responseBody.Data = map[string]interface{}{
		"domain":   "",
	}
	return &responseBody
}

// CreateUser 创建用户
func CreateUser(username string, password string) *ResponseBody {
	base64Pass := base64.StdEncoding.EncodeToString([]byte(password))
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	if username == "admin" {
		responseBody.Msg = "不能创建用户名为admin的用户!"
		return &responseBody
	}
	// create User
	sqlite := core.NewSqlite(DefaultDbPath)
	err := sqlite.CreateUserORM("", username, base64Pass, password)
	if err != nil {
		responseBody.Msg = "Failed to create"
	}
	return &responseBody
}

// UpdateUser 更新用户
func UpdateUser(id string, username string, password string) *ResponseBody {

	responseBody := ResponseBody{Msg: "success"}

	// 解码密码
	base64Pass, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		responseBody.Msg = "Base64解码失败: " + err.Error()
		return &responseBody
	}
	
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
	// 检查姓名是否重复
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
	// 检查密码是否重复
	if foundUser.Password != password {
		cond := struct{
			Password string
		}{
			Password: password,
		}
		if users, _ := sqlite.QueryUsersWhereORM(cond); len(users) != 0 {
			responseBody.Msg = "已存在密码为: " + password + " 的用户!"
			return &responseBody
		}
	}
	// 更新用户信息
	if err := sqlite.UpdateUserORM(id, username, password, string(base64Pass)); err != nil {
		responseBody.Msg = err.Error()
	}
	return &responseBody
}

// DelUser 删除用户
func DelUser(id string) *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
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
