package repository

import (
	"fmt"
	"gorm.io/gorm"
)

type Favorite struct {
	ID         int64
	UserID     int64
	BlogID     int64
	IsFavorite int32 `gorm:"default:0"`
}

type FavoriteDao struct {
}

// NewFavoriteDaoInstance 返回一个点赞表实体类的指针变量，可以方便调用该结构体的方法
func NewFavoriteDaoInstance() *FavoriteDao {
	return &FavoriteDao{}
}

// QueryFavoriteByUserIdAndBlogId 查询是否包含对应点赞关联记录
func (f *FavoriteDao) QueryFavoriteByUserIdAndBlogId(userId int64, blogId int64) (bool, error) {
	var count int64
	if err := db.Table("favorites").Where("user_id = ? and blog_id = ?", userId, blogId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// CreateFavorite 创建点赞关联记录
func (f *FavoriteDao) CreateFavorite(userId int64, blogId int64, actionType int32) error {
	blog := &Blog{}
	favorite := &Favorite{
		UserID:     userId,
		BlogID:     blogId,
		IsFavorite: actionType,
	}
	tx := db.Begin()
	//添加数据
	if err := tx.Select("user_id", "blog_id", "is_favorite").Create(favorite).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 点赞博文的点赞数+1
	if err := tx.Model(blog).Where("id = ? ", blogId).Update("favorite_count", gorm.Expr("favorite_count+ ?", 1)).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// DelBlogFavorite 删除博文点赞关联记录
func (f *FavoriteDao) DelBlogFavorite(blogId int64) error {
	favorite := &Favorite{}
	//添加数据
	if err := db.Where("blog_id = ?", blogId).Delete(favorite).Error; err != nil {
		return err
	}
	return nil
}

// UpdateFavorite 更新点赞关联记录
func (f *FavoriteDao) UpdateFavorite(userId int64, blogId int64, actionType int32) error {
	blog := &Blog{}
	favorite := &Favorite{
		UserID:     userId,
		BlogID:     blogId,
		IsFavorite: actionType,
	}
	fmt.Println("更新点赞标志位······")
	tx := db.Begin()
	//更新标志位
	if err := tx.Model(favorite).Where("user_id = ? and blog_id = ?", userId, blogId).Update("is_favorite", actionType).Error; err != nil {
		tx.Rollback()
		return err
	}
	if actionType == 1 {
		// 点赞博文的点赞数+1
		if err := tx.Model(blog).Where("id = ? ", blogId).Update("favorite_count", gorm.Expr("favorite_count+ ?", 1)).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		if err := tx.Model(blog).Where("id = ? ", blogId).Update("favorite_count", gorm.Expr("favorite_count- ?", 1)).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

// QueryBlogsIdByUserId 通过用户id查询查询该用户点赞的所有博文对应的博文id列表
func (f *FavoriteDao) QueryBlogsIdByUserId(userId int64) ([]int64, error) {
	var ids []int64
	fmt.Println("通过userId查询点赞博文列表的blogId")
	if err := db.Table("favorites").Select("blog_id").Where("user_id = ? and is_favorite = ?", userId, 1).Find(&ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

// QueryActionTypeByUserIdAndBlogId 通过用户id和博文id获取该用户对于这个博文是否点赞的状态码
func (f *FavoriteDao) QueryActionTypeByUserIdAndBlogId(userId, blogId int64) (int32, error) {
	var actionType int32
	fmt.Println("通过userId+blogId查询点赞状态")
	if err := db.Table("favorites").Select("is_favorite").Where("user_id = ? and blog_id = ?", userId, blogId).Find(&actionType).Error; err != nil {
		return 0, err
	}
	return actionType, nil
}

// QueryFavoriteCountByUserId 根据用户id查询用户点赞博文的数量
func (f *FavoriteDao) QueryFavoriteCountByUserId(userId int64) (int64, error) {
	var favoriteCount int64
	fmt.Println("通过userId查询点赞博文数量")
	if err := db.Model(&Favorite{}).Where("user_id = ? and is_favorite = ?", userId, 1).Count(&favoriteCount).Error; err != nil {
		return 0, err
	}
	return favoriteCount, nil
}
