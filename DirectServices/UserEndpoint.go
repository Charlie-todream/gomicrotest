package DirectServices


type UserRequest struct {
	Uid  int `json:"uid"`
	Method string
}

type UserResponse struct {
	Result string `json:"result"`
}
