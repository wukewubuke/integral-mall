package logic

import (
	"crypto/md5"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/yakaa/log4g"
	"my-integral-mall/common/baseerror"
	"my-integral-mall/common/rpcxclient/integralrpcmodel"
	"my-integral-mall/user/model"
	"strconv"
)

type (
	UserLogic struct {
		userModel *model.UserModel
		redisCache *redis.Client
		integralRpcModel *integralrpcmodel.IntegralRpcModel
	}

	RegisterRequest struct {
		Mobile string `json:"mobile" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	RegisterResponse struct {

	}

	LoginRequest struct {
		Mobile string `json:"mobile" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	LoginResponse struct {
		Authorization string
	}
)

var (
	ErrRecordExists = baseerror.NewBaseError("此手机号码已经存在")
	ErrUsernameOrPassword = baseerror.NewBaseError("用户名或密码错误")
)


func NewUserLogic(userModel *model.UserModel,redisCache *redis.Client,integralRpcModel *integralrpcmodel.IntegralRpcModel) *UserLogic {
	return &UserLogic{userModel:userModel, redisCache:redisCache, integralRpcModel: integralRpcModel}
}



func (l *UserLogic)Register(r *RegisterRequest)(*RegisterResponse,error){
	response := new(RegisterResponse)
	//判断手机号码是否存在
	b, err := l.userModel.ExistByMobile(r.Mobile)
	if err != nil {
		return  nil, err
	}
	if b {
		return nil, ErrRecordExists
	}


	user := &model.User{
		Mobile: r.Mobile,
		Password: fmt.Sprintf("%x", md5.Sum([]byte(r.Password))),
	}


	_, err = l.userModel.InsertTransaction(user, func(userId int64) error {
		if err := l.integralRpcModel.AddIntegral(userId, 1000); err != nil {
			log4g.ErrorFormat("add integral failed, error:%+v",err)
			return err
		}
		return nil
	})



	if err != nil {
		return nil, err
	}



	return response, nil
}


func (l *UserLogic)Login(r *LoginRequest)(*LoginResponse, error){
	response := new(LoginResponse)
	user, err := l.userModel.FindByMobile(r.Mobile)
	if err != nil {
		return nil, err
	}

	authorization := fmt.Sprintf("%x",md5.Sum([]byte(user.Mobile+strconv.Itoa(int(user.Id)))))
	l.redisCache.Set(authorization, user.Id, model.AuthorizationExpire)
	response.Authorization = authorization

	return  response, nil
	
}
