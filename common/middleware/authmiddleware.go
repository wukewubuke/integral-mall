package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"my-integral-mall/common/baseerror"
	"my-integral-mall/common/baseresponse"
	"strconv"
	"strings"
)


type (
	Authorization struct {
		redisCache *redis.Client
	}
)

var (
	ErrAuthorization = baseerror.NewBaseError("请先登录")
)

func NewAuthorization(redisCache *redis.Client)(*Authorization){
	return &Authorization{
		redisCache: redisCache,
	}
}

func (a *Authorization)Auth (ctx *gin.Context){
	authorization := ctx.GetHeader("Authorization")
	if strings.TrimSpace(authorization) == "" {
		baseresponse.HttpResponse(ctx, nil,ErrAuthorization)
		ctx.Abort()
		return
	}

	//从redis中取出用户id
	sc := a.redisCache.Get(authorization)
	userId ,_ := strconv.Atoi(sc.Val())
	ctx.Set("userId", userId)
	ctx.Next()
	return
 }
