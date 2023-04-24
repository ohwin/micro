package etcd

type Config struct {
	Endpoints []string `json:"endpoints" yaml:"endpoints"`
}

func (c Config) IsNil() bool {
	return len(c.Endpoints) == 0
}
