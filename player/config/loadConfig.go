package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func LoadAllConfig() (*AllConfigs, error) {
	//读取general_config
	filePath := fmt.Sprintf("./config/servers.json")
	buffer, err := readFile(filePath)
	if err != nil {
		return &AllConfigs{}, err
	}
	serverConfig := ServerConfig{}
	err = json.Unmarshal(buffer, &serverConfig)
	if err != nil {
		return &AllConfigs{}, err
	}
	allServerConfig := &AllConfigs{
		ServerConfig: serverConfig,
	}
	GetConfigManInstance().SetServiceConfig(allServerConfig)
	return allServerConfig, nil
}

func readFile(filePath string) ([]byte, error) {
	filePath, _ = filepath.Abs(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buffer, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}
