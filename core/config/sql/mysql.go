package sql

import "fmt"

type Config struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	DB       string `json:"db" yaml:"db"`
	Charset  string `json:"charset" yaml:"charset"`
}

func (c Config) IsNil() bool {
	return c == Config{}
}

func (c Config) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.DB,
		c.Charset,
	)
}
