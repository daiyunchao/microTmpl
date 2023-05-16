package configHandle

import "sync"

var once sync.Once
var instance *ConfigMan

type ConfigMan struct {
	AllConfigs *AllConfigs
}

func GetConfigManInstance() *ConfigMan {
	once.Do(func() {
		instance = &ConfigMan{}
	})
	return instance
}

func (config *ConfigMan) SetServiceConfig(allConfigs *AllConfigs) {
	config.AllConfigs = allConfigs
}

func (config *ConfigMan) GetGeneralConfig() GeneralConfig {
	return config.AllConfigs.GeneralConfig
}
