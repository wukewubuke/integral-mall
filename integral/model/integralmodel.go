package model

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
	"time"
)

type (
	Integral struct {
		Id         int64
		UserId     int64     `xorm:"int notnull 'user_id'"`
		Integral   int64     `xorm:"int notnull 'integral'"`
		CreateTime time.Time `xorm:"DateTime 'create_time'"`
	}

	IntegralModel struct {
		mysql      *xorm.Engine
		redisCache *redis.Client
		table      string
	}
)

func NewIntegralModel(mysql *xorm.Engine, redisCache *redis.Client, table string) *IntegralModel {
	return &IntegralModel{mysql: mysql, redisCache: redisCache, table: table}
}

func (m *IntegralModel) Insert(i *Integral) (int64, error) {
	return m.mysql.Insert(i)
}

func (m *IntegralModel) ExistsByUserId(userId int64) (bool, error) {
	return m.mysql.Exist(&Integral{UserId: userId})
}

func (m *IntegralModel) FindByUserId(userId int64) (*Integral, error) {
	integral := new(Integral)
	if _, err := m.mysql.Where("user_id = ?", userId).Get(integral); err != nil {
		return nil, err
	}
	return integral, nil
}

func (m *IntegralModel) FindById(id int64) (*Integral, error) {
	integral := new(Integral)
	if _, err := m.mysql.Where("id = ?", id).Get(integral); err != nil {
		return nil, err
	}
	return integral, nil
}

func (m *IntegralModel) UpdateIntegralByUserId(userId, integral int64) (*Integral, error) {
	query := "update " + m.table + " set integral = integral - ? where user_id = ?"
	if _, err := m.mysql.Exec(query, integral, userId); err != nil {
		return nil, err
	}
	return m.FindById(userId)
}

func (m *IntegralModel) InsertIntegralSql(userId, integral int64) string {
	return fmt.Sprintf("INSERT INTO %s (`user_id`, `integral`)values(%d,%d)",
		m.table, userId, integral)
}

func (m *IntegralModel) UpdateIntegralByUserIdSql(userId, integral int64) string {
	return fmt.Sprintf("update %s set integral = integral - %d where user_id = %d",
		m.table, integral, userId)
}

func (m *IntegralModel) ExecSql(sql string) error {
	if _, err := m.mysql.Exec(sql); err != nil {
		return err
	}
	return nil
}
