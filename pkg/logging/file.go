package logging

import (
	"fmt"
	"gin-blog/pkg/file"
	"gin-blog/pkg/setting"
	"os"
	"time"
)


/**
 * 获取日志文件保存路径
 * @return string
 * @date 2021-01-30 19:11:22
 * @author RedBo
 */
func getLogFilePath() string {
	return fmt.Sprintf("%s", setting.AppSetting.LogSavePath)
}

/**
 * 获取日志文件名称
 * @return string
 * @date 2021-01-30 19:12:01
 * @author RedBo
 */
func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		setting.AppSetting.LogSaveName,
		time.Now().Format(setting.AppSetting.TimeFormat),
		setting.AppSetting.LogFileExt,
	)
}

/**
 * 打开日志文件
 * @param filename
 * @param filePath
 * @return *os.File
 * @return error
 * @date 2021-01-30 19:15:23
 * @author RedBo
 */
func  openLogFile(filename,  filePath string) (*os.File, error)  {
	dir, err := os.Getwd()
	if err != nil {
		return nil , fmt.Errorf("os.Getwd err: %v", err)
	}
	src := dir + "/" + filePath
	perm := file.CheckNotExist(src)
	if perm == true {
		return nil , fmt.Errorf("file.CheckPermission Permission dennied src: %v", src)
	}

	err = file.IsNotExistMkdir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	f , err := file.Open(src + filename, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)

	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}
	return f, nil
}

func mkDir()  {
	dir, _ := os.Getwd()
	err := os.MkdirAll(dir + "/" + getLogFilePath(), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
