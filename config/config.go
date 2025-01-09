package config

import (
	"github.com/spf13/viper"
	"github.com/wonderivan/logger"
)

// Conf 全局变量
var Conf *Config

// Config 配置文件结构体, 与YAML文件对应
type Config struct {
	Server   string `yaml:"server"`
	MySQL    MySQL  `yaml:"mysql"`
	Redis    Redis  `yaml:"redis"`
	OSS      OSS    `yaml:"oss"`
	Settings map[string]interface{}
}

type MySQL struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	DB           string `yaml:"db"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	MaxIdleConns int    `yaml:"maxIdleConns"`
	MaxOpenConns int    `yaml:"maxOpenConns"`
	MaxLifeTime  int    `yaml:"maxLifeTime"`
}

type Redis struct {
	Host     string `yaml:"host"`
	DB       int    `yaml:"db"`
	Password string `yaml:"password"`
}

type OSS struct {
	Endpoint   string `yaml:"endpoint"`
	AccessKey  string `yaml:"accessKey"`
	SecretKey  string `yaml:"secretKey"`
	BucketName string `yaml:"bucketName"`
	SSL        bool   `yaml:"ssl"`
}

// InitConfig 配置文件初始化
func InitConfig() {

	v := viper.New()

	// 定义配置名称, 不需要带后缀
	v.SetConfigName("config")

	// 指定配置文件目录, 可以是相对路径
	v.AddConfigPath("config")

	// 指定配置文件类型
	v.SetConfigType("yaml")

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		logger.Error("配置文件初始化失败：" + err.Error())
		return
	}

	// 将配置文件反序列化成结构体
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		logger.Error("配置文件初始化失败：" + err.Error())
		return
	}

	// 将解析出来的配置赋值给全局变量
	Conf = &cfg
}
