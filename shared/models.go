package shared

import "encoding/json"

type User struct {
	UserName string   `json:"username"`
	Password string   `json:"password"`
	Cards    []string `json:"cards"`
	Deck     []string `json:"deck"`
}

type Request struct {
	Action string `json:"action"`
	Data   json.RawMessage `json:"data"`
}

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

/*type Card struct {
	NameCard string `json:"name"`
}
//Ideia futura, usar struct para as cartas e não string*/
