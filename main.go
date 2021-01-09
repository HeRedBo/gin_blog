package main

import (
	"fmt"
	"gin-blog/models"
	"gin-blog/pkg/logging"
	"gin-blog/pkg/setting"
	"gin-blog/routers"
	"log"
	"syscall"

	"github.com/fvbock/endless"
)

func main() {
	setting.Setup()
	models.Setup()
	logging.Setup()

	endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20
	endpoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)

	router := routers.InitRouter()
	server := endless.NewServer(endpoint, router)
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
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
