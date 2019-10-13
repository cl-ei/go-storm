package config

import (
	"github.com/wonderivan/logger"
	"gopkg.in/ini.v1"
	"os"
)

type Redis struct {
	host      string
	port      int
	password  string
	goStormDB int
}

type CloudFunc struct {
	agency   string
	acceptor string
	getUid   string
}

type Config struct {
	CDN_URL   string
	redis     Redis
	cloudFunc CloudFunc
}

var CONFIG Config

func init() {
	if logger.SetLogger("config/logger.json") != nil {
		os.Exit(-2)
		return
	}

	logger.Info("Loading configure file ...")
	configFileName := "/etc/madliar.settings.ini"

	_, err := os.Stat(configFileName)
	if err != nil {
		logger.Error("Error happened when loading configure file: ", err)
		os.Exit(-1)
		return
	}
	if os.IsNotExist(err) {
		logger.Error("Configure file not existed: ", err)
		os.Exit(-1)
		return
	}

	conf, err := ini.Load(configFileName)
	if err != nil {
		logger.Error("Error in loading config file: ", err)
		os.Exit(-1)
		return
	}

	CONFIG.CDN_URL = conf.Section("default").Key("CDN_URL").String()
	redisSec := conf.Section("redis")
	rPort, _ := redisSec.Key("port").Int()
	rgoStormDB, _ := redisSec.Key("go_storm_db").Int()
	CONFIG.redis = Redis{
		host:      redisSec.Key("host").String(),
		port:      rPort,
		password:  redisSec.Key("password").String(),
		goStormDB: rgoStormDB,
	}

	CloudFuncSec := conf.Section("cloud_function")
	CONFIG.cloudFunc = CloudFunc{
		agency:   CloudFuncSec.Key("url").String(),
		acceptor: CloudFuncSec.Key("acceptor").String(),
		getUid:   CloudFuncSec.Key("get_uid").String(),
	}

	logger.Info("Config File: ", CONFIG)
}
