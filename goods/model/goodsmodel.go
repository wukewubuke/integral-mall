package model

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
	"my-integral-mall/common/baseerror"
	"time"
)

type (
	/*
		1	id	int
		0	good_name	varchar	255
		0	price	int
		0	intro	text
		0	image	varchar	255
		0	store	int
		0	create_time	timestamp
	*/
	Goods struct {
		Id         int64
		GoodName   string    `xorm:"Varchar(255) 'good_name'"`
		Price      int       `xorm:"int 'price'"`
		Intro      string    `xorm:"text 'intro'"`
		Image      string    `xorm:"Varchar(255) 'image'"`
		Store      int       `xorm:"int 'store'"`
		CreateTime time.Time `xorm:"DateTime 'create_time'"`
	}

	GoodsModel struct {
		mysql      *xorm.Engine
		redisCache *redis.Client
		table      string
	}
)

const (
	goodsDefaultPageSize int = 10
)

var (
	ErrNotFound = baseerror.NewBaseError("没有找到相关记录")
)

func NewGoodsModel(mysql *xorm.Engine, redisCache *redis.Client, table string) *GoodsModel {
	return &GoodsModel{mysql: mysql, redisCache: redisCache, table: table}
}

func (m *GoodsModel) PageList(goodName string, page int) ([]*Goods, int64, error) {
	if page > 0 {
		page -= 1
	}
	goodsList := []*Goods(nil)
	count, err := m.mysql.Table(m.table).Where("name like ?", "%"+goodName+"%").Count(goodsList)
	if err != nil {
		return nil, 0, err
	}
	page = page * goodsDefaultPageSize
	if err := m.mysql.Table(m.table).Where("name like ? LIMIT ?, ?", "%"+goodName+"%", page, goodsDefaultPageSize).Find(&goodsList); err != nil {
		return nil, 0, err
	}
	return goodsList, count, nil
}


func (m *GoodsModel)FindById (id int64)(*Goods, error){
	goods := new(Goods)
	b, err := m.mysql.Table(m.table).Where("id = ?", id).Get(goods)
	if err != nil {
		return nil, err
	}

	if !b {
		return nil, ErrNotFound
	}
	return goods, nil
}

//事务下单后改变库存和用户积分
func (m *GoodsModel)TransactionChangeStore(id, num int64, userId int64,opts ...func(userId int64) error) error{
	_, err := m.mysql.Transaction(func(session *xorm.Session) (i interface{}, e error) {
		query := fmt.Sprintf("update %s set store = store - ? where id = ?", m.table)
		if  _, err := session.Exec(query, num, id ); err != nil {
			return nil, err
		}
		for _, opt := range opts {
			if err := opt(userId); err != nil {

				return nil, err
			}
		}

		return nil, nil
	})
	return  err
}
