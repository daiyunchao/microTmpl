package define

import "encoding/json"

type Request struct {
	Module string         `json:"module"` //模块名
	Method string         `json:"method"` //方法名
	Data   map[string]any `json:"data"`   //请求参数
	Ip     string         `json:"ip"`
}

func (request *Request) DecodeData(dist any) {
	bytes, err := json.Marshal(request.Data)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, dist)
	if err != nil {
		return
	}
}

type OriRequest struct {
	Data  string `json:"data"`
	Debug int    `json:"debug"`
}

type OriJsonRequest struct {
	Data map[string]any `json:"data"`
}
