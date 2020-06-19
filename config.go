package main

import (
	"io/ioutil"

	"github.com/valyala/fastjson"
)

// BotConfig конфиг бота однако
type BotConfig struct {
	Token           string
	GroupID         string
	Plugins         string
	ErrorDictionary string
	PlugConfigs     string
}

// эта функция парсит конфиг ("config.json")
func parseConfig(filename string) BotConfig {
	var cfg BotConfig
	dat, err := ioutil.ReadFile(filename)
	_check(err, true)
	var p fastjson.Parser
	v, err := p.Parse(string(dat))
	_check(err, true)
	cfg.Token = string(v.GetStringBytes("token"))
	cfg.GroupID = string(v.GetStringBytes("group_id"))
	cfg.Plugins = string(v.GetStringBytes("plugins"))
	cfg.ErrorDictionary = string(v.GetStringBytes("error_dictionary"))
	cfg.PlugConfigs = string(v.GetStringBytes("plugconfigs"))

	return cfg
}

func pushConfigs(p *Plugins, filename string) {
	dat, err := ioutil.ReadFile(filename)
	_check(err)

	var parser fastjson.Parser
	v, err := parser.Parse(string(dat))
	_check(err)

	for _, e := range *p {
		if e.Data.ConfigPos == "" {
			continue
		}

		if v.Exists(e.Data.ConfigPos) {
			e.Data.Config = v.Get(e.Data.ConfigPos)
		}
	}
}
