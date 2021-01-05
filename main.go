package main

import (
	"gin-blog/pkg/setting"
	"gin-blog/routers"
	"net/http"
	"fmt"
)

func main() {
	router := routers.InitRouter()
	s := &http.Server{
		Addr : fmt.Sprintf(":%d", setting.HttpPort),
		Handler: router,
		ReadTimeout: setting.ReadTimeout,
		WriteTimeout: setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()

}
