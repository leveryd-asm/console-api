package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Mysql struct {
		DataSource string
		Debug      bool
	}
	Cors struct {
		AllowOrigin string
	}
}
