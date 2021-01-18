package models

import "github.com/jinzhu/gorm"

type Article struct {
	Model

	TagID int `json:"tag_id" gorm:"index"`
	Tag Tag `json:tag`

	Title string `json:"title"`
	Desc string `json:"desc"`
	Content string `json:"content"`
	CreatedBy string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State int `json:"state"`
}

//func (tag *Article) BeforeCreate(scope *gorm.Scope) error {
//	scope.SetColumn("CreatedOn", time.Now().Unix())
//	return nil
//}
//
//func (tag *Article) BeforeUpdate(scope *gorm.Scope) error  {
//	scope.SetColumn("ModifiedOn",time.Now().Unix())
//	return nil
//}


func ExistArticleByID(id int) (bool, error) {
	var article Article
	err := db.Select("id").Where("id = ? ", id ).First(&article).Error

	if err !=nil && err != gorm.ErrRecordNotFound {
		return false ,err
	}

	if article.ID > 0 {
		return true, nil
	}
	return false, nil
}

/**
获取文章总数
 */
func GetArticleTotal(maps interface{}) (count int) {
	db.Model(&Article{}).Where(maps).Count(&count)
	return
}

/**
获取文章列表
 */
func GetArticles(pageNum int, pageSize int, maps interface{}) (articles []Article) {
	db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)
	return
}

/**
获取单个文章详情
 */
func GetArticle(id int) (*Article, error){
	//db.Where("id = ?", id).First(&article)
	//db.Model(&article).Related(&article.Tag)
	var article Article
	err := db.Where("id = ? AND deleted_on = ? ", id , 0 ). First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil ,err
	}

	err = db.Model(&article).Related(&article.Tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &article, nil
}

/**
文章添加
 */
func AddArticle(data map[string]interface{}) bool {
	db.Create(&Article {
		TagID : data["tag_id"].(int),
		Title : data["title"].(string),
		Desc : data["desc"].(string),
		Content: data["content"].(string),
		CreatedBy: data["created_by"].(string),
		State : data["state"].(int),
	})

	return true
}

/**
文章删除
 */
func DeleteArticle(id int) bool {
	db.Where("id = ?", id ).Delete(Article{})
	return true
}

func EditArticle(id int, data interface {}) bool {
	db.Model(&Article{}).Where("id = ?", id).Updates(data)

	return true
}

func CleanAllArticle() bool {
	db.Unscoped().Where("deleted_on != ? ", 0).Delete(&Article{})
	return true
}





