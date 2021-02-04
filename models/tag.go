package models

import "github.com/jinzhu/gorm"

type Tag struct {
	Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func GetTags(pageNum int, pageSize int, maps interface{}) ([]Tag, error) {
	var (
		tags []Tag
		err error
	)

	if pageSize > 0 && pageNum > 0 {
		err = db.Where(maps).Find(&tags).Offset(pageNum).Limit(pageSize).Error
	} else {
		err = db.Where(maps).Find(&tags).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return tags , nil

	//db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)
	//return

}

func GetTagTotal(maps interface{}) (count int) {
	db.Model(&Tag{}).Where(maps).Count(&count)

	return
}

func ExistTagByName(name string) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("name = ?", name).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false , err
	}
	if tag.ID > 0 {
		return true, nil
	}
	return false,nil
}

func AddTag(name string, state int, createdBy string) bool {
	db.Create(&Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	})
	return true
}

//func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
//	scope.SetColumn("CreatedOn", time.Now().Unix())
//	return nil
//}
//
//func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error  {
//	scope.SetColumn("ModifiedOn",time.Now().Unix())
//	return nil
//}

func ExistTagByID(ID int) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("id = ? ", ID).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}
	return false, nil
}

func EditTag(id int, data interface{}) bool {
	db.Model(&Tag{}).Where("id = ?", id).Update(data)
	return true
}

func DeleteTag(id int) bool {
	db.Where("id = ?", id).Delete(&Tag{})
	return true
}

func CleanAllTag() bool {
	db.Unscoped().Where("deleted_on != ? ", 0).Delete(&Tag{})
	return true
}
