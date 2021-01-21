package tag_service

import (
	"encoding/json"
	"gin-blog/models"
	"gin-blog/pkg/export"
	"gin-blog/pkg/gredis"
	"gin-blog/pkg/logging"
	"gin-blog/service/cache_service"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/tealeg/xlsx"
	"io"
	"strconv"
	"time"
)

type Tag struct {
	ID int
	Name string
	CreatedBy string
	ModifiedBy string
	State int

	PageNum int
	PageSize int
}

func (t *Tag) ExistByName() (bool, error) {
	//return models.ExistTagByName(t.name)
	return true, nil
}


func (t *Tag) GetALl() ([]models.Tag, error) {
	var (
		tags, cacheTags [] models.Tag
	)

	cache := cache_service.Tag{
		State:t.State,

		PageNum:t.PageNum,
		PageSize:t.PageSize,
	}

	key := cache.GetTagsKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err !=nil {
			logging.Error(err)
		} else {
			json.Unmarshal(data, &cacheTags)
			return cacheTags, nil
		}
	}

	tags, err := models.GetTags(t.PageNum, t.PageSize, t.getMaps())

	if err != nil {
		return nil ,err
	}

	gredis.Set(key, tags , 3600)
	return tags , nil
}

func (t *Tag) Export() (string , error) {
	tags , err := t.GetALl()
	if err != nil {
		return "",err
	}

	file := xlsx.NewFile()
	sheet , err := file.AddSheet("标签信息")
	if err !=nil {
		return "",err
	}
	titles := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	row := sheet.AddRow()

	var cell *xlsx.Cell
	for _, title := range titles {
		cell = row.AddCell()
		cell.Value = title
	}
	timeLayout := "2006-01-02 15:04:05"  //转化所需模板
	for _, v := range tags {
		values := []string{
			strconv.Itoa(v.ID),
			v.Name,
			v.CreatedBy,
			time.Unix(int64(v.CreatedOn), 0).Format(timeLayout),
			v.ModifiedBy,
			time.Unix(int64(v.ModifiedOn), 0).Format(timeLayout),
		}
		row = sheet.AddRow()
		for _, value := range values {
			cell = row.AddCell()
			cell.Value = value
		}
	}

	timeUnix := strconv.Itoa(int(time.Now().Unix()))
	filename := "tag_" + timeUnix + ".xlsx"
	fullPath := export.GetExcelFullPath() + filename
	logging.Info(fullPath)
	err = file.Save(fullPath)
	if err !=nil {
		return "",err
	}
	return filename,nil
}



func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string] interface{})
	maps["deleted_on"] = 0

	if t.Name != "" {
		maps["name"] = t.Name
	}

	if t.State >= 0 {
		maps["state"] = t.State
	}

	return maps
}

func (t *Tag) Import(r io.Reader) error {
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}

	rows := xlsx.GetRows("标签信息")
	for  irow, row := range rows {
		if irow > 0 {
			var data []string
			for _, cell := range row {
				data = append(data,cell)
			}
			models.AddTag(data[1], 1, data[2])
		}
	}
	return nil
}