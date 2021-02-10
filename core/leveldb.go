package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"

	"github.com/linuxing3/vpsman/util"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
)

// LoginInfo Only for Admin to store its information
type LoginInfo struct {
	Username string `json: "username"`
	Password string `json: "password"`
}

// GetValue 获取leveldb值
func GetValue(key string) (string, error) {
	fmt.Println(key)
	if runtime.GOOS == "windows" {
		return GetValueJSON(key)
	}
	// linux or macos
	var dbPath string = viper.GetString("main.db.leveldb.path")
	db, err := leveldb.OpenFile(dbPath, nil)
	defer db.Close()
	if err != nil {
		return "", err
	}
	result, err := db.Get([]byte(key), nil)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// SetValue 设置leveldb值
// admin的密码是保存在leveldb中
func SetValue(key string, value string) error {
	fmt.Println(key)
	fmt.Println(value)
	// windows
	if runtime.GOOS == "windows" {
		return SetValueJSON(key, value)
	}
	// linux and macos
	var dbPath string = viper.GetString("main.db.leveldb.path")
	db, err := leveldb.OpenFile(dbPath, nil)
	defer db.Close()
	if err != nil {
		return err
	}
	return db.Put([]byte(key), []byte(value), nil)
}

// DelValue 删除值
func DelValue(key string) error {
	if runtime.GOOS == "windows" {
		return SetValueJSON("", "")

	}
	var dbPath string = viper.GetString("main.db.leveldb.path")
	db, err := leveldb.OpenFile(dbPath, nil)
	defer db.Close()
	if err != nil {
		return err
	}
	return db.Delete([]byte(key), nil)
}

// GetValueJSON 从json文件读取
func GetValueJSON (key string) (string, error){
	var jsonPath string = viper.GetString("main.db.jsondb.path")
	util.EnsureFileExists(jsonPath)
	// windows
	loginInfo := LoginInfo{}
	data, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(data, &loginInfo); err != nil {
		return "", err
	}
	if loginInfo.Username == key {
		return loginInfo.Password, nil
	}
	return "", nil
	
}

// SetValueJSON 设置json文件
func SetValueJSON(key string, value string) error {
	var jsonPath string = viper.GetString("main.db.jsondb.path")
	util.EnsureFileExists(jsonPath)
	loginInfo := LoginInfo{
		Username: key,
		Password: value,
	}
	data, err := json.MarshalIndent(loginInfo, "", "    ");
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(jsonPath, data, 0644); err != nil {
		return err
	}
	return nil
}

