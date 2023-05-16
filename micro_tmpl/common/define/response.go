package define

type Response struct {
	Code  int    `json:"code"`
	Data  any    `json:"data"`
	Error string `json:"error"`
}

type NetResponse struct {
	Code  int    `json:"code"`
	Data  any    `json:"data"`
	Error string `json:"error"`
}
