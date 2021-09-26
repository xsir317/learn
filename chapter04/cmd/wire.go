//+build wireinject

package main

import (
	"go-learn/chapter04/configs"
	"go-learn/chapter04/dao"
	internal_dao "go-learn/chapter04/internal/dao"

	"github.com/google/wire"
)

func InitDao(ip string, port int) *dao.Dao {
	wire.Build(configs.NewConf, internal_dao.MakeFakeRedis)
	return *dao.Dao{}
}
