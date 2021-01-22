package qrcode

import (
	"gin-blog/pkg/file"
	"gin-blog/pkg/setting"
	"gin-blog/pkg/util"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/jpeg"
)

type Qrcode struct {
	URL string
	Width int
	Height int
	Ext string
	Level qr.ErrorCorrectionLevel
	Mode qr.Encoding
}

const (
	EXT_JPG = ".jpg"
)

func NewQrcode(url string, width, height int , level qr.ErrorCorrectionLevel, mode qr.Encoding)  *Qrcode {
	return &Qrcode{
		URL: url,
		Width: width,
		Height: height,
		Level: level,
		Mode: mode,
		Ext: EXT_JPG,
	}
}

// 获取二维码保存路径
func GetQrCodePath()  string {
	return setting.AppSetting.QrCodeSavePath
}

func GetQrcodeFullPath() string {
	return setting.AppSetting.RuntimeRootPath + setting.AppSetting.QrCodeSavePath
}

func GetQrCodeFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetQrCodePath() + name
}

func GetQrCodeFileName(value string) string {
	return util.EncodeMD5(value)
}

// GetQrCrCode get qr file ext
func (q *Qrcode) GetQrCodeExt() string {
	return q.Ext
}

func (q *Qrcode) Encode(path string) (string, string, error) {

	name := GetQrCodeFileName(q.URL) + q.GetQrCodeExt()
	src := path + name

	if file.CheckNotExist(src ) == true {
		code , err :=qr.Encode(q.URL, q.Level, q.Mode)
		if err != nil {
			return "", "", err
		}
		code ,err = barcode.Scale(code, q.Width, q.Height)
		if err != nil {
			return "", "", err
		}

		f, err := file.MustOpen(name, path)
		if err != nil {
			return  "", "", err
		}

		defer f.Close()

		err = jpeg.Encode(f, code, nil)
		if err != nil {
			return  "", "", err
		}
	}
	return name, path , nil

}









