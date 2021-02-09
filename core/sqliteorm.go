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
	var users []User
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
		if err := db.Find(&users, idsInt).Error; err != nil {
			return nil, err
		}
	} else {
		fmt.Println("Find all records:")
		if err := db.Where("id > ?", 0).Find(&users).Error; err != nil {
			return nil, err
		}
	}
	// 更改为指针数组
	for _, e := range users {
		userList = append(userList, &e)
	}
	return userList, nil
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

// QueryUsersWhereORM 根据指定多个id获取用户记录
func (s *Sqlite) QueryUsersWhereORM(cond interface{}) ([]*User, error) {
	var users []User
	var userList []*User
	db := s.Connect()
	fmt.Println("Find all records:")
	if err := db.Where(cond).Find(&users).Error; err != nil {
		return nil, err
	}
	// 更改为指针数组
	for _, e := range users {
		userList = append(userList, &e)
	}
	fmt.Println(users)
	return userList, nil
}
