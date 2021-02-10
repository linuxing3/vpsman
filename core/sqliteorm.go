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

// PageQueryUser 分页查询
type PageQueryUser struct {
	PageNum  int
	CurPage  int
	Total    int
	PageSize int
	DataList []*User
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
	encryPass := sha256.Sum224([]byte(originPass)) // %x => password
	if err := db.Create(&User{Username: username, Password: fmt.Sprintf("%x", encryPass), PasswordShow: base64Pass}).Error; err != nil {
		return err
	}
	return nil
}

// HasDuplicateUserORM 检查是否重复密码和用户名
func (s *Sqlite) HasDuplicateUserORM(username, password string) bool{
	var users []User
	encryPass := sha256.Sum224([]byte(password)) 
	db := s.Connect()
	db.Where("username = ? or password = ?", username, fmt.Sprintf("%x", encryPass)).Find(&users)
	if len(users) != 0 {
		return true
	}
	return false
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
	db.Save(&user)
	return nil
}

// UpdateUserCondORM 使用给定信息更新用户名和密码
func (s *Sqlite) UpdateUserCondORM(id string, data *User) error {
	// FIXED Do not use pointer, because User Struct Not initialized
	var user User
	db := s.Connect()
	db.First(&user, id)
	db.Model(&user).Updates(data)
	db.Save(&user)
	return nil
}

// DeleteUserORM 使用给定信息删除用户
func (s *Sqlite) DeleteUserORM(id string) error {
	// FIXED Do not use pointer, because User Struct Not initialized
	var user User
	db := s.Connect()
	if err := db.Delete(&user, id).Error; err != nil {
		return err
	}
	fmt.Println(&user)
	fmt.Println(user)
	return nil
}

// QueryUserORM 用id查询数据
func (s *Sqlite) QueryUserORM(id string) (*User, error) {
	var user User
	db := s.Connect()
	if err := db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// QueryUsersORM 根据指定多个id获取用户记录
func (s *Sqlite) QueryUsersORM(ids ...string) ([]*User, error) {
	// FIXED Use pointer
	var users []*User
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
		if err := db.Find(&users, idsInt).Error; err != nil {
			return nil, err
		}
	} else {
		fmt.Println("Find all records:")
		if err := db.Find(&users).Error; err != nil {
			return nil, err
		}
	}
	fmt.Println(users)
	return users, nil
}

// PageQueryUsersORM 分页查询用户信息
func (s *Sqlite)PageQueryUsersORM(curPage int, pageSize int) (*PageQueryUser, error) {
	var total int
	var users []*User
	offset := (curPage - 1) * pageSize
	db := s.Connect()
	if err := db.Offset(offset).Limit(pageSize).Find(&users).Error; err !=nil {
		return nil, err
	}
	return &PageQueryUser{
		CurPage:  curPage,
		PageSize: pageSize,
		Total:    total,
		DataList: users,
		PageNum:  (total + pageSize - 1) / pageSize,
	}, nil
}

// QueryUsersWithStructORM 根据Struct User获取用户记录
// When querying with struct, GORM will only query with non-zero fields, 
// that means if your field’s value is 0, '', false or other zero values, 
// it won’t be used to build query conditions
// Struct
// db.Where(&User{Name: "jinzhu", Age: 20}).First(&user)
// SELECT * FROM users WHERE name = "jinzhu" AND age = 20 ORDER BY id LIMIT 1;
// Slice of primary keys
// db.Where([]int64{20, 21, 22}).Find(&users)
// SELECT * FROM users WHERE id IN (20, 21, 22);
func (s *Sqlite) QueryUsersWithStructORM(cond *User) ([]*User, error) {
	var users []*User
	db := s.Connect()
	fmt.Println("Find all records with condition:")
	if err := db.Where(cond).Find(&users).Error; err != nil {
		return nil, err
	}
	fmt.Println(users)
	return users, nil
}

// QueryUsersWithInterface 根据map[string]interface{}获取用户记录
// Map
// db.Where(map[string]interface{}{"name": "jinzhu", "age": 20}).Find(&users)
// SELECT * FROM users WHERE name = "jinzhu" AND age = 20;
func (s *Sqlite) QueryUsersWithInterface(cond map[string]interface{}) ([]*User, error) {
	var users []*User
	db := s.Connect()
	fmt.Println("Find all records:")
	if err := db.Where(cond).Find(&users).Error; err != nil {
		return nil, err
	}
	fmt.Println(users)
	return users, nil
}
