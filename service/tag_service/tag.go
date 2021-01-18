package tag_service

import (
	"encoding/json"
	"gin-blog/models"
	"gin-blog/pkg/export"
	"gin-blog/pkg/gredis"
	"gin-blog/pkg/logging"
	"gin-blog/service/cache_service"
	"github.com/tealeg/xlsx"
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
	if err !=nil {
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

	for _, v := range tags {
		values := []string{
			strconv.Itoa(v.ID),
			v.Name,
			v.CreatedBy,
			strconv.Itoa(v.CreatedOn),
			v.ModifiedBy,
			strconv.Itoa(v.ModifiedOn),
		}
		row = sheet.AddRow()
		for _, value := range values {
			cell = row.AddCell()
			cell.Value = value
		}
	}


	time := strconv.Itoa(int(time.Now().Unix()))
	filename := "tag_" + time + ".xlsx"
	fullPath := export.GetExcelFullPath() + filename
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

