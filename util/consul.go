package util

import (
	"fmt"
	"github.com/google/uuid"
	consulapi "github.com/hashicorp/consul/api"
	"log"
	"strconv"
)
var ConsulClient *consulapi.Client
var ServiceId string
var ServiceName string
var ServicePort int

func init () {
	config := consulapi.DefaultConfig()
	config.Address="192.168.1.124:8500"
	client,err:= consulapi.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	ConsulClient = client
	ServiceId = "userservice" + uuid.New().String()
}

func SetServiceNameAndPort(name string,port int)  {
	ServiceName = name
	ServicePort = port
}
func RegService()  {
	config := consulapi.DefaultConfig()
	config.Address="192.168.1.124:8500"

	reg := consulapi.AgentServiceRegistration{}
	reg.ID =ServiceId
	reg.Name=ServiceName
	reg.Address="192.168.1.124"
	reg.Port=ServicePort
	reg.Tags=[]string{"primary"}
	check := consulapi.AgentServiceCheck{}
	check.Interval="5s"
	check.HTTP="http://192.168.1.124:"+strconv.Itoa(ServicePort)+"/health"

	reg.Check = &check


	fmt.Println("启动一个服务")
	ConsulClient.Agent().ServiceRegister(&reg)
}

func Unregservice(){
	ConsulClient.Agent().ServiceDeregister("userservice1")
}
