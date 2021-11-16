package main

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

// 直连
func main2() {
	tartget, _ := url.Parse("http://localhost:8080")

	// 直连client ,两个func 一个是如何请求 一个是响应我们怎么处理
	client := httptransport.NewClient("GET", tartget, DirectServices.GetUserInfoRequest, DirectServices.GetUserInfoResponse)
	// Endpoint返回调用远程HTTP端点的可用Go kit端点。
	getUserInfo := client.Endpoint()
	ctx := context.Background()

	res, err := getUserInfo(ctx, DirectServices.UserRequest{Uid: 102})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)

	}

	userInfo := res.(DirectServices.UserResponse)

	fmt.Println(userInfo.Result)

}

// 通过consul中心连接
func main() {
	// 第一步骤使用consul api 创建一个client 然后 go-kit的包帮我们封装好了一个专门的client
	{
		config := consulapi.DefaultConfig()
		config.Address = "localhost:8500"
		api_client, _ := consulapi.NewClient(config)
		client := consul.NewClient(api_client)

		var logger log.Logger
		{
			logger = log.NewLogfmtLogger(os.Stdout)
			var Tag = []string{"primary"}

			// 第二步骤创建一个consul的实例  true通过了检验才能得到  是服务名称而不是id
			instancer := consul.NewInstancer(client, logger, "userservice", Tag, true)
			{
				// factory定义了如何获得服务端的endpoint,这里的service_url是从consul中读取到的service的address我这里是192.168.1.124:8080
				factory := func(service_url string) (endpoint.Endpoint, io.Closer, error) {
					tgt, _ := url.Parse("http://" + service_url)
					return httptransport.NewClient("GET", tgt, DirectServices.GetUserInfoRequest, DirectServices.GetUserInfoResponse).Endpoint(), nil, nil
				}
				endpointer := sd.NewEndpointer(instancer, factory, logger)
				endpoints, err := endpointer.Endpoints()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println("服务有", len(endpoints), "条")

				//mylib := lb.NewRoundRobin(endpointer)  // 轮询的方式
				mylib := lb.NewRandom(endpointer, time.Now().UnixNano()) // 随机的方式
				for {

					getUserInfo, _ := mylib.Endpoint() // 写死获取第一个
					ctx := context.Background()        // 创建一个context上下文对象

					// 第四步:执行
					res, err := getUserInfo(ctx, DirectServices.UserRequest{
						Uid: 101,
					})

					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					// 第五步骤:断言，得到响应值

					userinfo := res.(DirectServices.UserResponse)
					fmt.Println(userinfo.Result)
					time.Sleep(time.Second * 3)
				}
			}
		}
	}
}
