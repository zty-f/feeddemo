package service

import (
	"errors"
	"feeddemo/common"
	"feeddemo/repository"
	"feeddemo/utils"
	"strconv"
)

var relationDaoInstance = repository.NewRelationDaoInstance()

type RelationService struct {
}

// NewRelationServiceInstance 返回一个关注服务类的指针变量，可以方便调用该结构体的方法
func NewRelationServiceInstance() *RelationService {
	return &RelationService{}
}

// DoRelationAction 关注
func (r *RelationService) DoRelationAction(userId, toUserId, actionType int64) error {
	//获取关注用户的所有博文列表
	feedsList, err := blogDaoInstance.QueryByOwner(toUserId)
	key := "following_feeds:" + strconv.FormatInt(userId, 10)
	if err != nil {
		return err
	}
	if actionType == 1 {
		// 关注
		err1 := relationDaoInstance.CreateRelation(userId, toUserId)
		if err1 != nil {
			return err1
		}
		for _, v := range feedsList {
			utils.ZAdd(key, v.UpdateTime.UnixNano(), v.Id)
		}
	} else {
		// 取消关注
		err2 := relationDaoInstance.DeleteRelation(userId, toUserId)
		if err2 != nil {
			return err2
		}
		for _, v := range feedsList {
			utils.ZRem(key, v.Id)
		}
	}
	return nil
}

// DoRelationFollowList 获取关注列表
func (r *RelationService) DoRelationFollowList(userId, loginUserId int64) ([]common.UserVo, error) {
	// 根据用户id查询该用户关注的所有用户的id
	ids, err := relationDaoInstance.QueryFollowIdsByUserId(userId)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, errors.New("还未关注任何用户，继续发现吧！")
	}
	// 根据用户的id查询用户信息
	users, err4 := userDaoInstance.QueryUsersByIds(ids)
	if err4 != nil {
		return nil, err4
	}
	userList := make([]common.UserVo, len(users))
	for i, _ := range users {
		favoriteCount, err1 := favoriteDaoInstance.QueryFavoriteCountByUserId(users[i].Id)
		if err1 != nil {
			return nil, err1
		}
		var totalFavorited = int64(0)
		if count, err2 := blogDaoInstance.QueryPublishCountByUserId(users[i].Id); err2 != nil {
			return nil, err2
		} else if count > 0 {
			totalFavorited, err = blogDaoInstance.QueryTotalFavoriteCountByUserId(users[i].Id)
		}
		if err != nil {
			return nil, err
		}
		isFollow, err3 := relationDaoInstance.QueryIsFollowByUserIdAndToUserId(loginUserId, users[i].Id)
		if err3 != nil {
			return nil, err3
		}
		userList[i] = common.UserVo{
			Id:             users[i].Id,
			UserName:       users[i].UserName,
			FollowCount:    users[i].FollowCount,
			FollowerCount:  users[i].FollowerCount,
			Avatar:         users[i].Avatar,
			IsFollow:       isFollow,
			TotalFavorited: totalFavorited,
			FavoriteCount:  favoriteCount,
		}
	}
	return userList, nil
}

// DoRelationFollowerList 获取粉丝列表
func (r *RelationService) DoRelationFollowerList(userId, loginUserId int64) ([]common.UserVo, error) {
	ids, err := relationDaoInstance.QueryFollowerIdsByUserId(userId)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, errors.New("还没有一个粉丝，有点可怜，继续创作吧！")
	}
	// 根据用户的id查询用户信息
	users, err1 := userDaoInstance.QueryUsersByIds(ids)
	if err1 != nil {
		return nil, err1
	}
	userList := make([]common.UserVo, len(users))
	for i, _ := range users {
		favoriteCount, err2 := favoriteDaoInstance.QueryFavoriteCountByUserId(users[i].Id)
		if err2 != nil {
			return nil, err2
		}
		var totalFavorited = int64(0)
		if count, err4 := blogDaoInstance.QueryPublishCountByUserId(users[i].Id); err4 != nil {
			return nil, err4
		} else if count > 0 {
			totalFavorited, err = blogDaoInstance.QueryTotalFavoriteCountByUserId(users[i].Id)
		}
		if err != nil {
			return nil, err
		}
		isFollow, err4 := relationDaoInstance.QueryIsFollowByUserIdAndToUserId(loginUserId, users[i].Id)
		if err4 != nil {
			return nil, err4
		}
		userList[i] = common.UserVo{
			Id:             users[i].Id,
			UserName:       users[i].UserName,
			FollowCount:    users[i].FollowCount,
			FollowerCount:  users[i].FollowerCount,
			Avatar:         users[i].Avatar,
			IsFollow:       isFollow,
			TotalFavorited: totalFavorited,
			FavoriteCount:  favoriteCount,
		}
	}
	return userList, nil
}
