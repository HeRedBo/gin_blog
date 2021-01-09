package logging

import (
	"fmt"
	"gin-blog/pkg/setting"
	"log"
	"os"
	"time"
)



func getLogFilePath() string {
	return fmt.Sprintf("%s", setting.AppSetting.LogSavePath)
}

func getLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("%s%s.%s",setting.AppSetting.LogSaveName, time.Now().Format(setting.AppSetting.TimeFormat), setting.AppSetting.LogFileExt)
	
	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}

func  openLogFile(filePath string) *os.File  {
	_, err := os.Stat(filePath)
	switch {
		case os.IsNotExist(err):
			mkDir()
		case os.IsPermission(err):
			log.Fatal("Permission :%v", err)
	}

	handle, err := os.OpenFile(filePath,os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Fail to OpenFile :%v", err)
	}

	return handle
}

func mkDir()  {
	dir, _ := os.Getwd()
	err := os.MkdirAll(dir + "/" + getLogFilePath(), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
