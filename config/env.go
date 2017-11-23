package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/golang/glog"
)

type LogConfigEntity struct {
	Dir   string `json:"dir"`
	Level string `json:"level"`
}

type ConfigEntity struct {
	DefaultServer string `json:"default_server"`
	DefaultPort   int    `json:"default_port"`

	GatewayServer string `json:"gateway_server"`
	GatewayPort   int    `json:"gateway_port"`

	UserServer string `json:"user_server"`
	UserPort   int    `json:"user_port"`

	EngineServer string `json:"engine_server"`
	EnginePort   int    `json:"engine_port"`

	LogCfg *LogConfigEntity `json:"log"`
	Roles  []string         `json:"roles"`
}

func (e *ConfigEntity) LoadFromFile(path string) error {
	if raw, err := ioutil.ReadFile(path); err != nil {
		glog.Errorf("read file(%s) error: err(%s)", path, err.Error())
		return err
	} else if err = json.Unmarshal(raw, e); err != nil {
		glog.Errorf("unmarshal error: err(%s)", err.Error())
		return err
	} else {
		glog.Infof("succeeded to load config from file(%s)", path)
		return nil
	}
}
