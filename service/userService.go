package service

import (
	"errors"
	"feeddemo/repository"
	"feeddemo/utils"
	"fmt"
)

const MaxUsernameLen = 32
const MaxPasswordLen = 32

var userDaoInstance = repository.NewUserDaoInstance()

type UserService struct {
}

// NewUserServiceInstance 返回一个用户服务类的指针变量，可以方便调用该结构体的方法
func NewUserServiceInstance() *UserService {
	return &UserService{}
}

// DoRegister 注册
func (u *UserService) DoRegister(userName, password string) error {
	fmt.Printf("用户正在注册：" + userName + ":" + password)
	uErr := checkUserName(userName)
	pErr := checkPassword(password)
	if uErr != nil {
		return uErr
	}
	if pErr != nil {
		return pErr
	}
	flag, err1 := userDaoInstance.QueryIsContainsUserName(userName)
	if err1 != nil {
		return err1
	}
	if flag {
		return errors.New("用户名已存在，请创建一个独一无二的name吧！")
	}
	password = utils.MD5(password)
	//调用Dao层
	_, err := userDaoInstance.CreateByNameAndPassword(userName, password)
	if err != nil {
		return err
	}
	return nil
}

//检查用户名
func checkUserName(userName string) error {
	if len(userName) > MaxUsernameLen {
		return errors.New("username is too long")
	}

	return nil
}

//检查密码
func checkPassword(passWord string) error {
	if len(passWord) > MaxPasswordLen {
		return errors.New("password is too long")
	}
	return nil
}
