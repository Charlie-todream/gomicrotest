**consul启动**

**服务编排的方式：**

sudo docker-compose -f docker/docker-compose.yml up

**用最简单的容器，启动**

docker run -d --name=cs -p 8500:8500 consul agent -server -bootsrap -ui -client 0.0.0.0

-ui 内置web 界面

-client -bootsrap 指定自己为leader

**查看服务列表**

http://192.168.1.124:8500/v1/agent/services

**手动注册服务提交**

curl
--request PUT
--data @p.json
http://127.0.0.1:8500/v1/agent/service/register?replace-existing-checks=true **// 注销服务** curl
--request PUT
http://127.0.0.1:8500/v1/agent/service/deregister/userservice1

**启动服务**

go run main.go --name userservice -p 8080
