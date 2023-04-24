package aliyun

import "github.com/ohwin/micro/core/config/aliyun/sms"

type Config struct {
	AccessKeyId     string     // AccessKeyId
	AccessKeySecret string     // AccessKeySecret
	Sms             sms.Config // 短信配置
}

func (c Config) IsNil() bool {
	return c == Config{}
}
