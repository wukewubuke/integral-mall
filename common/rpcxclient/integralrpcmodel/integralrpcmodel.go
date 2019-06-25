package integralrpcmodel

import (
	"context"
	"github.com/yakaa/grpcx"
	"github.com/yakaa/grpcx/config"
	"my-integral-mall/integral/protos"
)

type (
	IntegralRpcModel struct {
		client *grpcx.GrpcxClient
	}
	IntegralClientModel struct {
		UserId int64
		Integral int64
	}
)

func NewIntegralRpcModel(client *grpcx.GrpcxClient) *IntegralRpcModel {
	return &IntegralRpcModel{
		client: client,
	}
}

func (m *IntegralRpcModel) AddIntegral(userId, integral int64) error {
	conn, err := m.client.GetConnection()
	if err != nil {
		return err
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), config.GrpcxDialTimeout)
	defer cancelFunc()

	client := protos.NewIntegralRpcClient(conn)
	if _, err := client.AddIntegral(ctx, &protos.AddIntegralRequest{
		UserId:   int64(userId),
		Integral: int64(integral),
	}); err != nil {
		return err
	}

	return nil
}

func (m *IntegralRpcModel) ConsumerIntegral(userId, integral int64) error {
	conn, err := m.client.GetConnection()
	if err != nil {
		return err
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), config.GrpcxDialTimeout)
	defer cancelFunc()

	client := protos.NewIntegralRpcClient(conn)
	if _, err := client.ConsumerIntegral(ctx, &protos.ConsumerIntegralRequest{
		UserId:           userId,
		ConsumerIntegral: integral,
	}); err != nil {
		return err
	}

	return nil
}

func (m *IntegralRpcModel) FindOneByUserId(userId int64)(*IntegralClientModel,error){
	conn, err := m.client.GetConnection()
	if err != nil {
		return nil, err
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), config.GrpcxDialTimeout)
	defer cancelFunc()

	client := protos.NewIntegralRpcClient(conn)
	resp, err := client.FindOneByUserId(ctx, &protos.FindOneByUserIdRequest{
		UserId:           userId,
	})

	if  err != nil {
		return nil, err
	}

	return &IntegralClientModel{
		UserId: resp.UserId,
		Integral: resp.Integral,
	},nil
}
