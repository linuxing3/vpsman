package core

import (
	"crypto/sha256"
	"fmt"
	"strconv"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Sqlite 结构体
type Sqlite struct {
	Path     string `json:"path"`
}

// User Model
type User struct {
	gorm.Model
	ID           uint `gorm:"primarykey"`
	Username   string
	Password     string
	PasswordShow string
	Level        string
	Email        string
	Quota        int64
	Download     uint64
	Upload       uint64
	UseDays      uint
	ExpiryDate   string
}


// NewSqlite constructor
func NewSqlite (path string) *Sqlite {
	var defaultPath string = "./vpsman.db"
	if path == "" {
			path = defaultPath
	}
	return &Sqlite{
			Path: path,
	}
}

// Connect Connect Sqlite for UserModel
func (s *Sqlite)Connect() *gorm.DB {

	db, err := gorm.Open(sqlite.Open(s.Path), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{})
	return db
}

// CreateUserORM 使用给定信息创建成员
func (s *Sqlite) CreateUserORM(id string, username string, base64Pass string, originPass string) error {

	db := s.Connect()

	encryPass := sha256.Sum224([]byte(originPass))
	if err := db.Create(&User{Username: username, Password: fmt.Sprintf("%x", encryPass), PasswordShow: base64Pass}).Error; err != nil {
		return err
	}

	return nil
}

// UpdateUserORM 使用给定信息更新用户名和密码
func (s *Sqlite) UpdateUserORM(id string, username string, base64Pass string, originPass string) error {
	var user User
	db := s.Connect()

	encryPass := sha256.Sum224([]byte(originPass))
	if err := db.Where(&User{Username: username}).First(&user).Error; err != nil {
		return err
	}
	if err := db.Model(&user).Updates(&User{Password: fmt.Sprintf("%x", encryPass), PasswordShow: base64Pass}).Error; err != nil {
		return err
	}
	return nil

}

// DeleteUserORM 使用给定信息删除用户
func (s *Sqlite) DeleteUserORM(id string) error {
	var user User
	db := s.Connect()
	fmt.Println("Deleteing record:")
	fmt.Println(id)
	idInt, _ := strconv.Atoi(id)
	if err := db.Delete(&user, idInt).Error; err != nil {
		return err
	}
	return nil
}

// QueryUserORM 用id查询数据
func (s *Sqlite) QueryUserORM(id string) (*User, error) {
	var user User
	db := s.Connect()
	idInt, _ := strconv.Atoi(id)
	if err := db.Find(&user, idInt).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// QueryUsersORM 根据指定多个id获取用户记录
func (s *Sqlite) QueryUsersORM(ids ...string) ([]*User, error) {
	var user []User
	var userList []*User
	db := s.Connect()

	fmt.Println("Got records:")
	fmt.Println(len(ids))

	if len(ids) > 0 {
		fmt.Println("Find some records:")
		var idsInt []int
		for i, e := range ids {
			idInt, _ := strconv.Atoi(e)
			idsInt[i] = idInt
		}
		if err := db.Find(&user, idsInt).Error; err != nil {
			return nil, err
		}
	} else {
		fmt.Println("Find all records:")
		if err := db.Where("id > ?", 0).Find(&user).Error; err != nil {
			return nil, err
		}
	}
	// 更改为指针数组
	for _, e := range user {
		userList = append(userList, &e)
	}
	fmt.Println(user)
	return userList, nil
}

// QueryUsersWhereORM 根据指定多个id获取用户记录
func (s *Sqlite) QueryUsersWhereORM(cond interface{}) ([]*User, error) {
	var user []User
	var userList []*User
	db := s.Connect()
	fmt.Println("Find all records:")
	if err := db.Where(cond).Find(&user).Error; err != nil {
		return nil, err
	}
	// 更改为指针数组
	for _, e := range user {
		userList = append(userList, &e)
	}
	fmt.Println(user)
	return userList, nil
}
