package article_service

import (
	"encoding/json"
	"gin-blog/models"
	"gin-blog/pkg/gredis"
	"gin-blog/pkg/logging"
	"gin-blog/service/cache_service"
)

type Article struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	State         int
	CreatedBy     string
	ModifiedBy    string

	PageNum  int
	PageSize int
}

func (a *Article) Add() error  {
	article := map[string]interface{}{
		"tag_id": a.TagID,
		"title" : a.Title,
		"desc" : a.Desc,
		"content" : a.Content,
		"cover_image_url" : a.CoverImageUrl,
		"state" : a.State,
		"created_by" : a.CreatedBy,
	}

	if err := models.AddArticle(article); err != nil {
		return err
	}
	return nil
}

func (a *Article) ExistByID() (bool, error) {
	return models.ExistArticleByID(a.ID)
}

func (a *Article) Get() (*models.Article, error) {
	var cacheArticle *models.Article

	cache := cache_service.Article{ID: a.ID}

	key := cache.GetArticleKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else  {
			json.Unmarshal(data, &cacheArticle)
			return cacheArticle, nil
		}
	}

	article, err := models.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}

	gredis.Set(key, article, 3600 )
	return article, nil
}

func (a *Article) GetAll() ([]*models.Article , error){
	var (
		articles, cacheArticles []*models.Article
	)
	cache := cache_service.Article{
		ID: a.ID,
		State: a.State,

		PageNum: a.PageNum,
		PageSize: a.PageSize,
	}
	key := cache.GetArticleKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err !=nil {
			logging.Error(err)
		} else {
			json.Unmarshal(data, &cacheArticles)
			return cacheArticles,nil
		}
	}

	articles, err := models.GetArticles(a.PageNum, a.PageSize, a.getMaps())
	if err !=nil {
		return  nil, err
	}
	logging.Debug(articles)
	gredis.Set(key, articles, 3600)
	return articles,nil


}


func (a *Article) Count() (int , error){
	return models.GetArticleTotal(a.getMaps())
}


func (a *Article) getMaps() map[string]interface{}  {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	if a.State != -1 {
		maps["state"] = a.State
	}
	if a.TagID != -1 {
		maps["tag_id"] = a.TagID
	}
	return maps
}



