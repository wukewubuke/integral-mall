package controller

import (
	"github.com/gin-gonic/gin"
	"my-integral-mall/common/baseresponse"
	"my-integral-mall/goods/logic"
)

type (
	GoodsController struct{
		goodsLogic *logic.GoodsLogic
	}
)



func NewGoodsController(goodsLogic *logic.GoodsLogic) *GoodsController{
	return &GoodsController{
		goodsLogic: goodsLogic,
	}
}


func (c *GoodsController)GoodsSearch(ctx *gin.Context){
	r := new(logic.GoodSearchRequest)
	if err := ctx.ShouldBindJSON(r); err != nil {
		baseresponse.ParamError(ctx, err)
		return
	}
	res, err := c.goodsLogic.GoodsSearch(r)
	baseresponse.HttpResponse(ctx, res, err)

}


func (c *GoodsController)GoodsOrder(ctx *gin.Context){
	r := new(logic.GoodOrderRequest)
	if err := ctx.ShouldBindJSON(r); err != nil {
		baseresponse.ParamError(ctx, err)
		return
	}

	res, err := c.goodsLogic.GoodsOrder(r, ctx.GetInt("userId"))
	baseresponse.HttpResponse(ctx, res, err)
}






