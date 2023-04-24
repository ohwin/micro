package initialize

import (
	"github.com/ohwin/micro/core/config"
	"github.com/ohwin/micro/core/etcd"
	"github.com/ohwin/micro/core/log"
	"github.com/ohwin/micro/core/store/cache"
	"github.com/ohwin/micro/core/store/sql"
)

var inits []func()
var customInits []func()

func Init() {
	inits = []func(){
		config.Init,
		log.Init,
		etcd.Init,
		sql.Init,
		cache.Init,
	}

	inits = append(inits, customInits...)

	for _, f := range inits {
		f()
	}
}

// Add 添加自定义初始化函数
func Add(fu ...func()) {
	for _, f := range fu {
		customInits = append(customInits, f)
	}
}
