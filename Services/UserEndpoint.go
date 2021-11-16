package Services

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/time/rate"
	"net/http"
	"service.gomicro.test/util"
	"strconv"
)

type UserRequest struct {
	Uid    int `json:"uid"`
	Method string
}

type UserResponse struct {
	Result string `json:"result"`
}

func GenUserEndpoint(userService IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserRequest)
		result := "nothing"
		if r.Method == "GET" {
			result = userService.GetName(r.Uid) + strconv.Itoa(util.ServicePort)
		} else if r.Method == "DELETE" {
			err := userService.DelUser(r.Uid)
			if err != nil {
				result = err.Error()
			}
		} else {
			result = fmt.Sprintf("userId为%d的用户删除成功！", r.Uid)
		}

		return UserResponse{Result: result}, nil
	}
}

// 加入限流功能
func RateLimit(limit *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !limit.Allow() {
				return nil, util.NewMyError(429, "too many request")

			}
			return next(ctx, request)
		}
	}
}

func MyErrorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	contentType, body := "text/plain;charset=utf-8", []byte(err.Error())
	w.Header().Set("Content-type", contentType) // 设置请求头

	if myerr, ok := err.(*util.MyError); ok {
		w.WriteHeader(myerr.Code)
		w.Write(body)
	} else {
		w.WriteHeader(500)
		w.Write(body)
	}

	w.WriteHeader(500)
	w.Write(body)
}
