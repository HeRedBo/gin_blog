package article_service

import (
	"gin-blog/pkg/file"
	"gin-blog/pkg/qrcode"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
)



type ArticlePoster struct {
	posterName string
	*Article
	Qr *qrcode.Qrcode
}


func NewArticlePoster(posterName string , article *Article, qr *qrcode.Qrcode) *ArticlePoster {
	return &ArticlePoster{
		posterName :posterName,
		Article : article,
		Qr: qr,
	}
}

func GetPosterFlag() string {
	return "poster"
}

func (a *ArticlePoster) ChekMergedImage(path string) bool {
	if file.CheckNotExist(path + a.posterName) == true {
		return false
	}
	return true
}

func (a *ArticlePoster) OpenMergedImage(path string) (*os.File, error) {
	f, err := file.MustOpen(a.posterName, path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

type Rect struct {
	Name string
	X0 int
	Y0 int
	X1 int
	Y1 int
}

type Pt struct {
	X int
	Y int
}

type ArticlePosterBg struct {
	Name string
	*ArticlePoster
	*Rect
	*Pt
}


func NewArticlePosterBg(name string ,ap *ArticlePoster, rect *Rect, pt *Pt) *ArticlePosterBg {
	return &ArticlePosterBg{
		Name : name,
		ArticlePoster:ap,
		Rect: rect,
		Pt:pt,
	}
}

func (a *ArticlePosterBg) Generate() (string, string , error) {
	fullPath := qrcode.GetQrcodeFullPath()
	fileName, path , err := a.Qr.Encode(fullPath)
	if err != nil {
		return "", "", err
	}

	if ! a.ChekMergedImage(path) {
		mergedF , err := a.OpenMergedImage(path)
		if err != nil {
			return "", "", err
		}
		defer mergedF.Close()
		
		bgF, err := file.MustOpen(a.Name, path)
		if err != nil {
			return "","", err
		}
		defer bgF.Close()

		qrF, err := file.MustOpen(fileName, path)
		if err != nil {
			return "","", err
		}
		defer qrF.Close()


		bgImage, err := jpeg.Decode(qrF)
		if err != nil {
			return "", "", err
		}

		qrImage , err := jpeg.Decode(qrF)
		if err != nil {
			return "","", err
		}
		jpg := image.NewRGBA(image.Rect(a.X0, a.Rect.Y0,a.X1, a.Rect.Y1))
		draw.Draw(jpg, jpg.Bounds(),bgImage,bgImage.Bounds().Min,draw.Over)

		draw.Draw(jpg, jpg.Bounds(), qrImage, qrImage.Bounds().Min.Sub(image.Pt(a.Pt.X, a.Pt.Y)), draw.Over)

		jpeg.Encode(mergedF, jpg, nil)

	}
	return fileName, path, nil
}



type DrawText struct {
	JPG draw.Image
	Merged *os.File

	Title string
	X0 int
	Y0 int
	Size0 float64

	SubTitle string
	X1 int
	Y1 int
	Size1 float64

}


