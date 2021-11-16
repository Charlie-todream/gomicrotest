package util

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	consulapi "github.com/hashicorp/consul/api"
	"io"
	"net/url"
	"os"
	"service.gomicro.test/DirectServices"
	"time"
)

func GetUser() (string, error) {
	//第一步创建client
	{
		config := consulapi.DefaultConfig()            //初始化consul的配置
		config.Address = "localhost:8500"              //consul的地址
		api_client, err := consulapi.NewClient(config) //根据consul的配置初始化client
		if err != nil {
			return "", err
		}
		client := consul.NewClient(api_client) //根据client创建实例

		var logger log.Logger
		{
			logger = log.NewLogfmtLogger(os.Stdout)
			var Tag = []string{"primary"}
			instancer := consul.NewInstancer(client, logger, "userservice", Tag, true) //最后的true表示只有通过健康检查的服务才能被得到
			{
				factory := func(service_url string) (endpoint.Endpoint, io.Closer, error) { //factory定义了如何获得服务端的endpoint,这里的service_url是从consul中读取到的service的address我这里是192.168.3.14:8000
					tart, _ := url.Parse("http://" + service_url)                                                                                 //server ip +8080真实服务的地址
					return httptransport.NewClient("GET", tart, DirectServices.GetUserInfoRequest, DirectServices.GetUserInfoResponse).Endpoint(), nil, nil //我再GetUserInfo_Request里面定义了访问哪一个api把url拼接成了http://192.168.3.14:8000/v1/user/{uid}的形式
				}
				endpointer := sd.NewEndpointer(instancer, factory, logger)
				endpoints, err := endpointer.Endpoints() //获取所有的服务端当前server的所有endpoint函数
				if err != nil {
					return "", err
				}
				fmt.Println("服务有", len(endpoints), "条")

				mylb := lb.NewRandom(endpointer, time.Now().UnixNano()) //使用go-kit自带的轮询

				for {
					getUserInfo, err := mylb.Endpoint() //写死获取第一个
					ctx := context.Background()         //第三步：创建一个context上下文对象

					//第四步：执行
					res, err := getUserInfo(ctx, DirectServices.UserRequest{Uid: 101})
					if err != nil {
						return "", err
					}
					//第五步：断言，得到响应值
					userinfo := res.(DirectServices.UserResponse)
					return userinfo.Result, nil
				}

			}
		}
	}
}
