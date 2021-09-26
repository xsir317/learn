package dao

import (
	"go-learn/chapter04/configs"
	"go-learn/chapter04/dao"
)

type FakeRedis struct {
	conf *configs.Conf
}

func MakeFakeRedis(c *configs.Conf) *dao.Dao {
	return &FakeRedis{conf: c}
}

func (redis *FakeRedis) Get(key string) (val string, err error) {
	return "Fake value", nil
}

func (redis *FakeRedis) Set(key string, val string) error {
	return nil
}
