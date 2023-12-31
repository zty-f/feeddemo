package repository

import (
	"fmt"
	"gorm.io/gorm"
)

type Relation struct {
	ID          int64
	UserID      int64 `gorm:"default:0"`
	FollowingID int64 `gorm:"default:0"`
}

type RelationDao struct {
}

// NewRelationDaoInstance 返回一个关注表实体类的指针变量，可以方便调用该结构体的方法
func NewRelationDaoInstance() *RelationDao {
	return &RelationDao{}
}

// CreateRelation 新增关注记录信息
func (r *RelationDao) CreateRelation(userId, toUserId int64) error {
	user := &User{}
	relation := &Relation{
		UserID:      userId,
		FollowingID: toUserId,
	}
	tx := db.Begin()
	if err := tx.Select("user_id", "following_id").Create(&relation).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 对应用户粉丝数+1
	if err := tx.Model(user).Where("id = ? ", toUserId).Update("follower_count", gorm.Expr("follower_count+ ?", 1)).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 对应用户关注数量+1
	if err := tx.Model(user).Where("id = ? ", userId).Update("follow_count", gorm.Expr("follow_count+ ?", 1)).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// DeleteRelation 删除关注记录信息
func (r *RelationDao) DeleteRelation(userId, toUserId int64) error {
	user := &User{}
	relation := &Relation{
		UserID:      userId,
		FollowingID: toUserId,
	}
	tx := db.Begin()
	if err := tx.Where("user_id = ? and following_id = ?", userId, toUserId).Delete(relation).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 对应用户粉丝数-1
	if err := tx.Model(user).Where("id = ? ", toUserId).Update("follower_count", gorm.Expr("follower_count- ?", 1)).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 对应用户关注数量-1
	if err := tx.Model(user).Where("id = ? ", userId).Update("follow_count", gorm.Expr("follow_count- ?", 1)).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// QueryIsFollowByUserIdAndToUserId 通过登录用户id和博文发布者id获取该登录用户是否关注博文所有者
func (r *RelationDao) QueryIsFollowByUserIdAndToUserId(userId, toUserId int64) (bool, error) {
	var count int64
	fmt.Println("通过userId+toUserId查询关注状态")
	if err := db.Table("relations").Select("count(1)").Where("user_id = ? and following_id = ?", userId, toUserId).Limit(1).Count(&count).Error; err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

// QueryFollowIdsByUserId 通过用户id查询该用户关注的所有用户的id
func (r *RelationDao) QueryFollowIdsByUserId(userId int64) ([]int64, error) {
	var ids []int64
	fmt.Println("通过用户id查询该用户关注的所有用户的id")
	if err := db.Table("relations").Select("following_id").Where("user_id = ?", userId).Find(&ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

// QueryFollowerIdsByUserId 通过用户id查询该用户所有粉丝的用户id
func (r *RelationDao) QueryFollowerIdsByUserId(userId int64) ([]int64, error) {
	var ids []int64
	fmt.Println("通过用户id查询该用户所有粉丝的id")
	if err := db.Table("relations").Select("user_id").Where("following_id = ?", userId).Find(&ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}
