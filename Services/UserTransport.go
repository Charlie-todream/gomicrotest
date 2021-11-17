package Services

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func DecodeUserRequest(c context.Context, r *http.Request) (interface{},error) {

	// r.URL.Query().Get("uid") != ""

	vars := mux.Vars(r)
	if uid,ok := vars["uid"];ok {
			uid,_:= strconv.Atoi(uid)
		return UserRequest{Uid: uid, Method: r.Method, Token: r.URL.Query().Get("token")}, nil
			//请求必须携带token过来，如果找不到这里返回空字符串，因为request访问的先后顺序是先DecodeUserRequest，
		// 再EncodeUserResponse再到我们的EndPoint，所以这里就已经给我们的request结构体存入了Token，那么我们EndPoint里面的request类型断言成UserRequest结构体实例后里面就有Token了
		}
		return nil,errors.New("参数错误")
	}

func EncodeUserResponse(ctx context.Context,w http.ResponseWriter,response interface{}) error {
	w.Header().Set("Content-type","application/json")
	 return  json.NewEncoder(w).Encode(response)
}