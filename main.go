package main

import (
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	. "service.gomicro.test/Services"
)

func main() {
	user := UserService{}
	endp:=GenUserEndpoint(user)
	serverHandler := httptransport.NewServer(endp,DecodeUserRequest,EncodeUserResponse)
	router :=mux.NewRouter()
	{
		router.Methods("GET","DELETE").Path(`'/user/{uid:\d+}'`).Handler(serverHandler)
		router.Methods("GET").Path("/health").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-type","application/json")
			writer.Write([]byte(`{"status":"ok"}`))
		})
	}
	http.ListenAndServe(":8080",router)
}
