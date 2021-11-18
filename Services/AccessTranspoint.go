package Services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

// gjson无法获取 post 数据
func DecodeAccessRequest(c context.Context, r *http.Request) (interface{}, error) {
	body, _ := ioutil.ReadAll(r.Body)

	result := gjson.Parse(string(body)) //第三方库解析json

	if result.IsObject() { //如果是json就返回true
		username := result.Get("username")
		userpass := result.Get("userpass")
		fmt.Println(userpass)
		fmt.Println(username)
		return AccessRequest{Username: username.String(), Userpass: userpass.String(), Method: r.Method}, nil
	}
	return nil, errors.New("用户名密码错误")

}
func EncodeAccessResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response) //返回一个bool值判断response是否可以正确的转化为json，不能则抛出异常，返回给调用方
}
