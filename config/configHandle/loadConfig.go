package configHandle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func LoadAllConfig(env string) (*AllConfigs, error) {
	//读取general_config
	filePath := fmt.Sprintf("../config/%s/general_config.json", env)
	buffer, err := readFile(filePath)
	if err != nil {
		return &AllConfigs{}, err
	}
	generalConfig := GeneralConfig{
		NotVerifyRoute: make([]string, 0),
	}
	err = json.Unmarshal(buffer, &generalConfig)
	if err != nil {
		return &AllConfigs{}, err
	}
	allServerConfig := &AllConfigs{
		GeneralConfig: generalConfig,
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
