package sms

import "time"

// Config 阿里云短信配置
type Config struct {
	SignName     string        // 短信签名
	TemplateCode string        // 短信模板
	Expire       time.Duration // 短信过期时间
}
