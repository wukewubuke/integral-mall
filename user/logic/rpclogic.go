package logic

import (
	"context"
	"my-integral-mall/user/model"
	"my-integral-mall/user/protos"
)

type (
	UserRpcServerLogic struct {
		UserModel *model.UserModel
	}

)


func NewUserRpcServerLogic(userModel *model.UserModel)(*UserRpcServerLogic){
	return &UserRpcServerLogic{
		UserModel: userModel,
	}
}


func (l *UserRpcServerLogic)FindByMobile(_ context.Context,r  *protos.FindByMobileRequest) (*protos.UserResponse, error){
	user, err := l.UserModel.FindByMobile(r.Mobile)
	if err != nil {
		return nil, err
	}
	resp := &protos.UserResponse{
		Id: user.Id,
		Name: user.Name,
		Mobile: user.Mobile,
	}
	return resp, err
}
func (l *UserRpcServerLogic)FindById(_ context.Context,r *protos.FindByIdRequest) (*protos.UserResponse, error){
	user, err := l.UserModel.FindById(r.Id)
	if err != nil {
		return nil, err
	}
	resp := &protos.UserResponse{
		Id: int64(user.Id),
		Name: user.Name,
		Mobile: user.Mobile,
	}
	return resp, err
}
