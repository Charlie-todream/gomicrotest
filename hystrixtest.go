package main

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"math/rand"
	"sync"
	"time"
)

type Product struct {
	ID    int
	Title string
	Price int
}

func getProduct() (Product, error) {
	r := rand.Intn(10)
	if r < 6 { //模拟api卡顿和超时效果
		time.Sleep(time.Second * 4)
	}
	return Product{
		ID:    101,
		Title: "Golang从入门到精通",
		Price: 12,
	}, nil
}

func RecProduct() (Product, error) {
	return Product{
		ID:    999,
		Title: "推荐商品",
		Price: 120,
	}, nil

}

func main() {

	rand.Seed(time.Now().UnixNano())
	configA := hystrix.CommandConfig{ //创建一个hystrix的config
		Timeout:               3000, //command运行超过3秒就会报超时错误
		MaxConcurrentRequests: 5,    //控制最大并发数为5，如果超过5会调用我们传入的回调函数降级
		RequestVolumeThreshold: 5,  // 在一个统计窗口没处理的请求量达到阈值，才会进行熔断与否的判断
		ErrorPercentThreshold: 20, //  在一个 %20的处理失败 处理熔断服务
		SleepWindow: int(time.Second * 10), // 熔断后多久尝试是否恢复
	}
	hystrix.ConfigureCommand("get_prod", configA) //hystrix绑定command
	c,_,_:=hystrix.GetCircuit("get_prod") // 熔断指针 ，bool表示是否能取到 error
	resultChan := make(chan Product, 1)

	wg := sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		go (func() {
			wg.Add(1)
			defer wg.Done()
			// Go为异步
			errs := hystrix.Do("get_prod", func() error { //使用hystrix来讲我们的操作封装成command,hystrix返回值是一个chan error
				p, _ := getProduct() //这里会随机延迟0-4秒
				resultChan <- p
				return nil //这里返回的error在回调中可以获取到，也就是下面的e变量
			}, func(e error) error {
				fmt.Println(e)
				rcp, err := RecProduct() //推荐商品,如果这里的err不是nil,那么就会忘errs中写入这个err，下面的select就可以监控到
				resultChan <- rcp
				return err
			})
			if errs!=nil {
				fmt.Println(errs)
			}else {
				select {
				case getProd := <-resultChan:
					fmt.Println(getProd)
				}
			}

			fmt.Println(c.IsOpen())
			fmt.Println(c.AllowRequest())
		})()

	}
	wg.Wait()

	time.Sleep(time.Second * 1)
}