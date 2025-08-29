package shared

type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type Request struct{
	Action string `json:"action"`
	Data interface{} `json:"data"`
}

type Response struct {
	Status string `json:"status"`
	Message string `json:"message,omitempty"`
	Data interface{} `json:"data,omitempty"`
}