package main

import (
	"fmt"
	"gin-blog/pkg/setting"
	"gin-blog/routers"
	"github.com/fvbock/endless"
	"log"
	"syscall"
)

func main() {
	endless.DefaultReadTimeOut = setting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20
	endpoint := fmt.Sprintf(":%d", setting.HttpPort)

	router := routers.InitRouter()

	server := endless.NewServer(endpoint, router)
	server.BeforeBegin = func(add string ) {
		log.Printf("Actual pid is %d", syscall.Getpid() )
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}
	
	//server := &http.Server{
	//	Addr : fmt.Sprintf(":%d", setting.HttpPort),
	//	Handler: router,
	//	ReadTimeout: setting.ReadTimeout,
	//	WriteTimeout: setting.WriteTimeout,
	//	MaxHeaderBytes: 1 << 20,
	//}
	//server.ListenAndServe()

}
