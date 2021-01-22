package file

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
)


/**
 * 获取文件大小
 * @param f
 * @return int
 * @return error
 * @date 2021-01-10 22:45:13
 * @author RedBo
 */
func GetSize(f multipart.File) (int, error) {
	content ,err := ioutil.ReadAll(f)

	return len(content), err
}

/**
获取文件后缀
 */
func GetExt(fileName string) string {
	return path.Ext(fileName)
}

/**
检查问价是否存在
 */
func CheckNotExist(src string) bool {
	_, err := os.Stat(src)

	return os.IsNotExist(err)
}

/**
 * 检查文件权限
 * @param src string
 * @return error
 */
func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	return os.IsPermission(err)
}

/**
 * 如果不存在就新建文件夹
 * @param src
 * @return error
 * @date 2021-01-10 23:02:37
 */
func IsNotExistMkdir(src string) error {
	if noExist := CheckNotExist(src) ; noExist == false {
		if err := Mkdir(src) ; err != nil {
			return err
		}
	}
	return nil
}


/**
 * 新建文件夹
 * @param src
 * @return error
 */
func Mkdir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 打开文件
 * @param name
 * @param flag
 * @param perm
 * @return *os.File
 * @return error
 */
func Open(name string, flag int , perm os.FileMode) (*os.File , error ) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func MustOpen(fileName, filePath string) (*os.File,  error) {
	dir , err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + "/" + filePath
	perm := CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file CheckPermission permission denied err: %s", src)
	}

	err = IsNotExistMkdir(src)
	if err != nil {
		return nil, fmt.Errorf("file IsNotExistMkdir src: %s , err : %v", src,err)
	}

	f, err := Open(src + fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("file OpenFile src:  %v", err)
	}
	return f,nil
}