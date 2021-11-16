package main

import (
	"flag"
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"os"
	"os/signal"
	"service.gomicro.test/Services"
	"service.gomicro.test/util"
	"strconv"
	"syscall"
)

func main() {

	name := flag.String("name", "", "服务名称")
	port := flag.Int("p", 0, "服务端口")
	flag.Parse()

	if *name == "" {
		log.Fatal("请指定服务名")
	}

	if *port == 0 {
		log.Fatal("请指定端口")
	}
	util.SetServiceNameAndPort(*name, *port)

	user := Services.UserService{}

	// 限流调用
	limit := rate.NewLimiter(1, 5)
	endp := Services.RateLimit(limit)(Services.GenUserEndpoint(user))

	options := []httptransport.ServerOption{
		//  生产ServerOption 切片,传入我们自定义的错误处理
		httptransport.ServerErrorEncoder(Services.MyErrorEncoder),
	}

	serverHandler := httptransport.NewServer(endp, Services.DecodeUserRequest, Services.EncodeUserResponse, options...)
	router := mux.NewRouter()

	router.Methods("GET", "DELETE").Path(`/user/{uid:\d+}`).Handler(serverHandler)

	router.Methods("GET").Path("/health").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-type", "application/json")
		writer.Write([]byte(`{"status":"ok"}`))
	})

	errChan := make(chan error)
	go func() {
		util.RegService() // 注册consul服务
		err := http.ListenAndServe(":"+strconv.Itoa(*port), router)
		if err != nil {
			log.Println(err)
		}
		errChan <- err
	}()
	// 优雅的关闭信号监听
	go func() {
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-sigChan)
	}()

	getErr := <-errChan
	util.Unregservice()
	log.Println(getErr)

}
