package Services

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"golang.org/x/time/rate"
	"net/http"
	"service.gomicro.test/util"
	"strconv"
)

type UserRequest struct {
	Uid    int `json:"uid"`
	Method string
	Token  string
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

// 加入限流功能中间件
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

//token验证中间件
func CheckTokenMiddleware() endpoint.Middleware { //Middleware type Middleware func(Endpoint) Endpoint
	return func(next endpoint.Endpoint) endpoint.Endpoint { //Endpoint type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			r := request.(UserRequest) //通过类型断言获取请求结构体
			uc := UserClaim{}
			//下面的r.Token是在代码DecodeUserRequest那里封装进去的
			getToken, err := jwt.ParseWithClaims(r.Token, &uc, func(token *jwt.Token) (i interface{}, e error) {
				return []byte(secKey), err
			})
			fmt.Println(err, 123)
			if getToken != nil && getToken.Valid { //验证通过
				newCtx := context.WithValue(ctx, "LoginUser", getToken.Claims.(*UserClaim).Uname)
				return next(newCtx, request)
			} else {
				return nil, util.NewMyError(403, "error token")
			}

			//logger.Log("method", r.Method, "event", "get user", "userid", r.Uid)

		}
	}
}

// 日志中间件
func UserServiceLogMiddleware(logger log.Logger) endpoint.Middleware  {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			r := request.(UserRequest)
			logger.Log("method",r.Method,"event","get user",r.Uid)
			return next(ctx,request)
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
