package conf

import (
	//	"strings"
	"syscommon"
	"time"

	"github.com/pelletier/go-toml"
)

type Configure struct {

	//System config
	ServicePort  int
	VCodeTimeOut int64

	//DB config
	DBType           string
	ConnectionString string
	LogMode          bool

	//	//Redis
	//	RedisProtocol    string
	//	RedisAddress     string
	//	RedisMaxPoolSize int

	//Redis-cluster
	RedisClusterAddresses    []string
	RedisClusterConnTimeout  time.Duration
	RedisClusterReadTimeout  time.Duration
	RedisClusterWriteTimeout time.Duration
	RedisClusterKeepAlive    int
	RedisClusterAliveTime    time.Duration
}

var configure *Configure

func GetConfigure() *Configure {

	return configure
}

func LoadConfigure(fileName string) error {

	config, err := toml.LoadFile(fileName)

	if err != nil {
		return err
	}

	conf := Configure{}
	conf.ServicePort = int(config.Get("sysconf.ServicePort").(int64))
	conf.VCodeTimeOut = config.Get("sysconf.VCodeTimeOut").(int64)

	//	conf.DBType = config.Get("db.DBType").(string)
	//	conf.ConnectionString = config.Get("db.ConnectionString").(string)
	//	conf.LogMode = config.Get("db.LogMode").(bool)

	//	conf.RedisAddress = config.Get("redis.Address").(string)
	//	conf.RedisMaxPoolSize = int(config.Get("redis.MaxPoolSize").(int64))
	//	conf.RedisProtocol = config.Get("redis.Protocol").(string)
	conf.RedisClusterAddresses = syscommon.InterfacesToStrings(config.Get("redis-cluster.Addresses").([]interface{}))
	conf.RedisClusterConnTimeout = time.Duration(config.Get("redis-cluster.ConnTimeout").(int64)) * time.Millisecond
	conf.RedisClusterReadTimeout = time.Duration(config.Get("redis-cluster.ReadTimeout").(int64)) * time.Millisecond
	conf.RedisClusterWriteTimeout = time.Duration(config.Get("redis-cluster.WriteTimeout").(int64)) * time.Millisecond
	conf.RedisClusterKeepAlive = int(config.Get("redis-cluster.KeepAlive").(int64))
	conf.RedisClusterAliveTime = time.Duration(config.Get("redis-cluster.AliveTime").(int64)) * time.Millisecond

	configure = &conf

	return nil

}
