package gatewayDefine

type ReqLoginData struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type ReqRegisterData struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
