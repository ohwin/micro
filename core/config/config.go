package config

import (
	"fmt"
	"github.com/ohwin/micro/core/config/aliyun"
	"github.com/ohwin/micro/core/config/cache"
	"github.com/ohwin/micro/core/config/etcd"
	"github.com/ohwin/micro/core/config/sql"
	"github.com/spf13/viper"
)

var App Config
var configPath string

type Config struct {
	Env    string        `json:"env" yaml:"env"` // dev prod test
	Host   string        `json:"host" yaml:"host"`
	Port   int           `json:"port" yaml:"port"`
	Mysql  sql.Config    `json:"mysql" yaml:"mysql"`
	Cache  cache.Config  `json:"cache" yaml:"cache"`
	Aliyun aliyun.Config `json:"aliyun" yaml:"aliyun"`
	Etcd   etcd.Config   `json:"etcd" yaml:"etcd"`
}

type NilConfig interface {
	IsNil() bool
}

func Addr() string {
	return fmt.Sprintf("%s:%d", App.Host, App.Port)
}

func Load(path string) {
	configPath = path
}

func Parse(path string) {
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&App)
	if err != nil {
		panic(err)
	}
}

func Init() {
	Parse(configPath)
}
