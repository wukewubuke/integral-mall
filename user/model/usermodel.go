package model

import (
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
	"time"
)

type (

	User struct {
		Id         int64
		Name       string    `xorm:"Varchar(20) notnull 'name'"`
		Mobile     string    `xorm:"Varchar(25) notnull unique 'mobile'"`
		Password   string    `xorm:"Varchar(32) notnull 'password'"`
		CreateTime time.Time `xorm:"DateTime 'create_time'"`
	}

	UserModel struct {
		mysql      *xorm.Engine
		redisCache *redis.Client
		table      string
	}
)

const (
	AuthorizationExpire  = 604800 * time.Second //7 * 24 *3600
)


func NewUserModel(mysql *xorm.Engine, redisCache *redis.Client, table string) *UserModel {
	return &UserModel{mysql: mysql, redisCache: redisCache, table: table}
}

func (m *UserModel) Insert(u *User) (int64, error) {
	return m.mysql.Insert(u)
}


func (m *UserModel)InsertTransaction(u *User, opts ...func(userId int64) error)(*User, error){
	_, err := m.mysql.Transaction(func(session *xorm.Session) (i interface{}, e error) {
		if  _, err := session.Insert(u); err != nil {
			return nil, err
		}
		for _, opt := range opts {
			if err := opt(u.Id); err != nil {

				return nil, err
			}
		}

		return u, nil
	})
	return u, err
}

//判断手机号码是否存在
func (m *UserModel) ExistByMobile(mobile string) (bool, error) {
	return m.mysql.Exist(&User{Mobile: mobile})
}


//根据手机号码取用户信息
func (m *UserModel) FindByMobile(mobile string) (*User, error) {
	user := new(User)
	if _, err :=  m.mysql.Where("mobile = ?", mobile).Get(user); err != nil {
		return  nil, err
	}
	return user , nil
}

//根据Id取用户信息
func (m *UserModel) FindById(id int64) (*User, error) {
	user := new(User)
	if _, err :=  m.mysql.Where("id = ?", id).Get(user); err != nil {
		return  nil, err
	}
	return user , nil
}
