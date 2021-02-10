package util

import (
	"fmt"
	"io/ioutil"
)

// EnsureFileExists 创建空文件
func EnsureFileExists(path string) {
	fmt.Printf("文件路径是 %s", path)
	if !IsExists(path) {
		// mainConf := viper.Get(key)
		// mainConfMap := map[string]interface{}{key: mainConf}
		// mainData, _ := json.Marshal(mainConfMap)
		ioutil.WriteFile(path, []byte{}, 0644)
	}
}