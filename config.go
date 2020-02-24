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
	SQL             struct {
		Username string
		Password string
		DBName   string
	}
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
	// mysql := v.GetObject("mysql")
	// cfg.SQL.Username = string(mysql.Get("username").GetStringBytes())
	// cfg.SQL.Password = string(mysql.Get("password").GetStringBytes())
	// cfg.SQL.DBName = string(mysql.Get("db").GetStringBytes())
	return cfg
}
