package main

import (
	"log"
	"service.gomicro.test/util"
)

func main() {
	err := util.GenRSAPubAndPri(1024, "./pem") //1024是长度，长度越长安全性越高，但是性能也就越差
	if err != nil {
		log.Fatal(err)
	}
	//执行完生成公钥和私钥，公钥给别人私钥给自己
}