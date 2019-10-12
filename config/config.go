package config

import (
	"fmt"
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
	fmt.Println("Loading configure file ...")
	configFileName := "/etc/madliar.settings.ini"

	_, err := os.Stat(configFileName)
	if err != nil {
		fmt.Println("Error happened when loading configure file: ", err)
		os.Exit(-1)
		return
	}
	if os.IsNotExist(err) {
		fmt.Println("Configure file not existed: ", err)
		os.Exit(-1)
		return
	}

	conf, err := ini.Load(configFileName)
	if err != nil {
		fmt.Print("Error in loading config file: ", err)
		os.Exit(-1)
		return
	}

	CONFIG.CDN_URL = conf.Section("default").Key("CDN_URL").String()
	r := Redis{}
	redisSec := conf.Section("redis")
	r.host = redisSec.Key("host").String()
	r.port, _ = redisSec.Key("port").Int()
	r.password = redisSec.Key("password").String()
	r.goStormDB, _ = redisSec.Key("go_storm_db").Int()
	CONFIG.redis = r

	c := CloudFunc{}
	CloudFuncSec := conf.Section("cloud_function")
	c.agency = CloudFuncSec.Key("url").String()
	c.acceptor = CloudFuncSec.Key("acceptor").String()
	c.getUid = CloudFuncSec.Key("get_uid").String()
	CONFIG.cloudFunc = c
}
