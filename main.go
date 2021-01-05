package main

import (
	"gin-blog/pkg/setting"
	"gin-blog/routers"
	"github.com/fvbock/endless"
	"net/http"
	"fmt"
)

func main() {
	endless.DefaultReadTimeOut = setting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20
	//endpoint := fmt.Sprintf("%d", setting.HttpPort)

	router := routers.InitRouter()
	server := &http.Server{
		Addr : fmt.Sprintf(":%d", setting.HttpPort),
		Handler: router,
		ReadTimeout: setting.ReadTimeout,
		WriteTimeout: setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()

}
