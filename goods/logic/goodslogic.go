package logic

import (
	"my-integral-mall/common/baseerror"
	"my-integral-mall/common/rpcxclient/integralrpcmodel"
	"my-integral-mall/goods/model"
)

type (
	GoodsLogic struct {
		goodsModel       *model.GoodsModel
		integralRpcModel *integralrpcmodel.IntegralRpcModel
	}

	GoodSearchRequest struct {
		Name string `json:"name"`
		Page int `json:"page"`
	}
	GoodSearchResponse struct {
		Total int64 `json:"total"`
		GoodsList []*GoodView `json:"goods_list"`
	}
	GoodView struct {
		Id    int64    `json:"id"`
		Name  string `json:"name"`
		Image string `json:"image"`
		Intro string `json:"intro"`
		Price int    `json:"price"`
		Store int    `json:"store"`
	}

	GoodOrderRequest struct {
		Id  int64 `json:"id"  binding:"required"`
		Num int64 `json:"num" binding:"required"`
	}


	GoodOrderResponse struct {

	}
)

var (
	ErrStoreOver = baseerror.NewBaseError("商品库存不足")
	ErrIntegralOver = baseerror.NewBaseError("积分不足")
)



func NewGoodsLogic(goodsModel *model.GoodsModel, integralModel *integralrpcmodel.IntegralRpcModel) *GoodsLogic {
	return &GoodsLogic{
		goodsModel:       goodsModel,
		integralRpcModel: integralModel,
	}
}


func (l *GoodsLogic) GoodsSearch(r *GoodSearchRequest) (*GoodSearchResponse, error) {

	goodsList,count,  err := l.goodsModel.PageList(r.Name, r.Page)
	if err != nil {
		return  nil, err
	}
	response := &GoodSearchResponse{ Total: count}
	for _, goods := range goodsList {
		response.GoodsList = append(response.GoodsList, &GoodView{
			Id: goods.Id,
			Name: goods.GoodName,
			Image: goods.Image,
			Intro: goods.Intro,
			Price: goods.Price,
			Store: goods.Store,
		})
	}
	return  response, nil
}

//下单
func (l *GoodsLogic) GoodsOrder(r *GoodOrderRequest, userId int) (*GoodOrderResponse, error) {
	//先查询商品是否存在
	goods, err := l.goodsModel.FindById(r.Id)
	if err != nil {
		return nil, err
	}
	if goods.Store <= 0 {
		return nil, ErrStoreOver
	}
	integral, err := l.integralRpcModel.FindOneByUserId(int64(userId))
	if err != nil {
		return nil, err
	}
	if integral.Integral < int64(goods.Price) {
		return nil, ErrIntegralOver
	}

	if err := l.goodsModel.TransactionChangeStore(r.Id, r.Num, int64(userId),  func(userId int64) error {
		//积分修改
		return l.integralRpcModel.ConsumerIntegral(userId, int64(goods.Price))

	}); err != nil {
		return nil, err
	}


	return  &GoodOrderResponse{

	},nil

}
