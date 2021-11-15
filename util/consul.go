package util

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"log"
)
var ConsulClient *consulapi.Client
func init () {
	config := consulapi.DefaultConfig()
	config.Address="192.168.1.124:8500"
	client,err:= consulapi.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	ConsulClient = client
}
func RegService()  {
	config := consulapi.DefaultConfig()
	config.Address="192.168.1.124:8500"

	reg := consulapi.AgentServiceRegistration{}
	reg.ID ="userservice1"
	reg.Name="userservice"
	reg.Address="192.168.1.124"
	reg.Port=8080
	reg.Tags=[]string{"primary"}
	check := consulapi.AgentServiceCheck{}
	check.Interval="5s"
	check.HTTP="http://192.168.1.124:8080/health"

	reg.Check = &check


	fmt.Println("启动一个服务")
	ConsulClient.Agent().ServiceRegister(&reg)
}

func Unregservice(){
	ConsulClient.Agent().ServiceDeregister("userservice1")
}
