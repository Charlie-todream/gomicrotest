package DirectServices

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// 是对pathinfo 部分进行拼接 path部分
func GetUserInfoRequest(_ context.Context,request *http.Request, r interface{}) error {
	user_request := r.(UserRequest)
	request.URL.Path += "/user/" + strconv.Itoa(user_request.Uid)
	return nil
}

func GetUserInfoResponse(_ context.Context, res *http.Response) (response interface{}, err error)  {
	if res.StatusCode > 400 {
		return nil,errors.New("No data")
	}
	var user_response UserResponse
	err = json.NewDecoder(res.Body).Decode(&user_response)

	if err != nil {
		return nil,err
	}
	return user_response,err
}
