package cache

type Config struct {
	Host     string
	Password string
	DB       int
}

func (c Config) IsNil() bool {
	return c == Config{}
}
