package main

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"log"
	"service.gomicro.test/util"
	"time"
)



func main() {
	configA := hystrix.CommandConfig{
		Timeout:                2000,
		MaxConcurrentRequests:  5,
		RequestVolumeThreshold: 3,
		SleepWindow:            int(time.Second * 10),
		ErrorPercentThreshold:  20,
	}

	hystrix.ConfigureCommand("getuser", configA)
	err := hystrix.Do("getuser", func() error {
		res, err := util.GetUser() //调用方法
		fmt.Println(res)
		return err
	}, func(e error) error {
		fmt.Println("降级用户")
		return e
	})
	if err != nil {
		log.Fatal(err)
	}
}
