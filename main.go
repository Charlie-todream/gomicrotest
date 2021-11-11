package main

import (
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
	. "service.gomicro.test/Services"
)

func main() {
	user := UserService{}
	endp:=GenUserEndpoint(user)
	serverHandler := httptransport.NewServer(endp,DecodeUserRequest,EncodeUserResponse)
	http.ListenAndServe(":8080",serverHandler)
}
