package config

type ServerConfig struct {
	EtcdUrl        string `json:"etcdUrl"`
	EtcdServerName string `json:"etcdServerName"`
}

type AllConfigs struct {
	ServerConfig
}
