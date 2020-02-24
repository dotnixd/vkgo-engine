package main

import (
	"fmt"
	"log"
	"path/filepath"
	"plugin"
	"strings"

	"./plugins/basic"
)

// Plugin описывает плагин для облегчения работы
type Plugin struct {
	Data basic.Plugin
	Sym  *plugin.Plugin
}

// Plugins массив плагинов
type Plugins []Plugin

// Load загружает плагины в ОЗУ
func (p *Plugins) Load(path string) {
	files, err := filepath.Glob(path + "/*/*.so")
	_check(err, true)

	log.Println(":: Найдены плагины: ")

	for _, filename := range files {
		file, err := plugin.Open(filename)
		_check(err, true)

		sym, err := file.Lookup("Init")
		_check(err)
		var plugin Plugin
		var data basic.Plugin

		sym.(func(*basic.Plugin))(&data)
		plugin.Data = data
		plugin.Sym = file

		fmt.Println("-", data.Cmds, "-", data.Desc)
		*p = append(*p, plugin)
	}
}

// Sym ищет подходящий плагин для обработки сообщения
func (p *Plugins) Sym(cmd string) (*Plugin, bool) {
	for _, sym := range *p {
		for _, name := range sym.Data.Cmds {
			if name == cmd {
				return &sym, true
			}
		}
	}

	return &Plugin{}, false
}

// Launch запускает обработку сообщения
func (p *Plugin) Launch() basic.APIError {
	run, err := p.Sym.Lookup("Run")
	_check(err)
	p.Data.Values.Args = strings.Fields(string(p.Data.Data.GetStringBytes("text")))
	p.Data.Values.PeerID = p.Data.Data.GetInt64("peer_id")
	p.Data.Values.FromID = p.Data.Data.GetInt64("from_id")
	return run.(func(*basic.Plugin) basic.APIError)(&p.Data)
}
