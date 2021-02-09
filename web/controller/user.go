package controller

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
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
	userWithName := &core.User{
		Username: findUser,
	}
	if findUser != "" {
		if users, _ = sqlite.QueryUsersWithStructORM(userWithName); len(users) > 0 {
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
	if sqlite.HasDuplicateUserORM(username, password) {
		responseBody.Msg = "已存在这个用户名或密码的用户!"
		return &responseBody
	}
	base64Pass := base64.StdEncoding.EncodeToString([]byte(password)) // passwordShow
	if err := sqlite.CreateUserORM("", username, base64Pass, password); err != nil {
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

	foundUser, _ := sqlite.QueryUserORM(id);
	if foundUser == nil {
		responseBody.Msg = "不存在这个用户id, 无法修改!"
		return &responseBody
	}

	if !sqlite.HasDuplicateUserORM(username, password) {
		responseBody.Msg = "不存在这个用户名，无法修改!"
		return &responseBody
	}
	
	encryPass := sha256.Sum224([]byte(password)) // %x => password
	base64Pass := base64.StdEncoding.EncodeToString([]byte(password)) // passwordShow
	data := core.User{
		Username: username,
		Password: fmt.Sprintf("%x", encryPass),
		PasswordShow: base64Pass,
	}
	if err := sqlite.UpdateUserCondORM(id, &data); err != nil {
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
