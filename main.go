package main

import (
	"log"

	"./errorapi"
	"./utils"
	"./vkapi"
)

// обработка исключений
func _check(err error, die ...bool) {
	if err != nil {
		if len(die) == 1 {
			log.Fatalln(err)
		} else {
			log.Println("Некритичная ошибка:", err)
		}
	}
}

func main() {

	log.Println(":: Запускаю бота...")
	var vk vkapi.VK
	cfg := parseConfig("./config.json")
	vk.Init(cfg.Token, cfg.GroupID)

	var PSys Plugins
	PSys.Load(cfg.Plugins)

	var ErrorHandler errorapi.ErrorHandler
	ErrorHandler.Init(&vk, cfg.ErrorDictionary)

	for {
		Handle(&vk, &PSys, &ErrorHandler)
	}
}

// Handle крутится в бесконечном цикле и получает новые сообщения
func Handle(vk *vkapi.VK, p *Plugins, e *errorapi.ErrorHandler) {
	upd := vk.ListenLongPoll()
	for _, v := range upd.Updates {
		if v.Type == "message_new" {
			sym, found := p.Sym(utils.ExtractCmd(string(v.Object.GetStringBytes("text"))))
			if found {
				sym.Data.VK = vk
				sym.Data.Data = v.Object
				go func() {
					err := sym.Launch()
					e.Handle(&sym.Data, &err)
				}()
			}
		}
	}
}
