package common

import "time"

//公共返回对象包

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type UserVo struct {
	Id             int64  `json:"id"`
	UserName       string `json:"username"`
	Avatar         string `json:"avatar"`
	FollowCount    int64  `json:"follow_count,omitempty"`
	FollowerCount  int64  `json:"follower_count,omitempty"`
	TotalFavorited int64  `json:"total_favorited,omitempty"` //获赞数
	FavoriteCount  int64  `json:"favorite_count,omitempty"`  //点赞数
	IsFollow       bool   `json:"is_follow"`
}

type BlogVo struct {
	Id            int64     `json:"id"`
	Author        UserVo    `json:"author"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	FavoriteCount int64     `json:"favorite_count"`
	CommentCount  int64     `json:"comment_count"`
	Top           string    `json:"top"`
	IsFavorite    bool      `json:"is_favorite"`
	CreateTime    time.Time `json:"create_time,omitempty"`
	UpdateTime    time.Time `json:"update_time,omitempty"`
}

type CommentVo struct {
	Id         int64  `json:"id,omitempty"`
	User       UserVo `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}
