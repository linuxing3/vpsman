package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	// "github.com/linuxing/vpsman/util"
)

// ServerConfig 服务器结构体定义
type ServerConfig struct {}

// Load 加载服务端配置文件
func Load(path string) *ServerConfig {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("加载配置文件失败")
		fmt.Println(err)
		return nil
	}
	config := ServerConfig{}
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Println("json写入xray失败")
		fmt.Println(err)
		return nil
	}
	return &config
}

// Save 保存服务端配置文件
func Save(config *ServerConfig, path string) bool {
	fmt.Println("保存xray服务端配置文件")
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		fmt.Println("xray服务端配置文件MarshalIndent失败")
		fmt.Println(err)
		return false
	}
	if err = ioutil.WriteFile(path, data, 0644); err != nil {
		fmt.Println("保存xray服务端配置文件失败")
		fmt.Println(err)
		return false
	}
	return true
}

// WriteTls 写tls配置
func WriteTls(cert, key, domain string) bool {
	config := Load("")
	return Save(config, "")
}

// WriteDomain 写域名
func WriteDomain(domain string) bool {
	config := Load("")
	// config.Inbounds[0].StreamSettings.SNI = domain
	return Save(config, "")
}

// WriteLogLevel 写日志等级
func WriteLogLevel(level string) bool {
	config := Load("")
	return Save(config, "")
}
