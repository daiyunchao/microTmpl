package configHandle

type GeneralConfig struct {
	AppId       string `json:"appId"`
	AppSecret   string `json:"appSecret"`
	IsEncode    bool   `json:"isEncode"`
	RedisConfig struct {
	} `json:"redisConfig"`
	MongodbConfig struct {
		ConnAddress string `json:"connAddress"`
		DBName      string `json:"dbName"`
	} `json:"mongodbConfig"`
	EtcdConfig struct {
		EtcdUrl        string `json:"etcdUrl"`
		EtcdServerName string `json:"etcdServerName"`
	}
	NotVerifyRoute []string `json:"notVerifyRoute"`
	AdminConfig    struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}
	BILog struct {
		TDSwitch bool `json:"tdswitch"`
		TD       struct {
			TDType     string `json:"tdtype"`
			URL        string `json:"url"`
			BIFilePath string `json:"bifilepath"`
			AppID      string `json:"appid"`
			BatchSize  int    `json:"batchsize"`
		} `json:"td"`
	} `json:"biLog"`
}

type AllConfigs struct {
	GeneralConfig
}
