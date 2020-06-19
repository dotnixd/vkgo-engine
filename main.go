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

	pushConfigs(&PSys, cfg.PlugConfigs)

	var ErrorHandler errorapi.ErrorHandler
	ErrorHandler.Init(&vk, cfg.ErrorDictionary)

	helpmsg := PSys.Help()

	for {
		Handle(&vk, &PSys, &ErrorHandler, helpmsg)
	}
}

// Handle крутится в бесконечном цикле и получает новые сообщения
func Handle(vk *vkapi.VK, p *Plugins, e *errorapi.ErrorHandler, helpmsg string) {
	upd := vk.ListenLongPoll()
	for _, v := range upd.Updates {
		if v.Type == "message_new" {
			cmd := utils.ExtractCmd(string(v.Object.GetStringBytes("text")))

			sym, found := p.Sym(cmd)
			if found {
				sym.Data.VK = vk
				sym.Data.Data = v.Object
				go func() {
					err := sym.Launch()
					e.Handle(&sym.Data, &err)
				}()
			}

			if cmd == "хелп" || cmd == "помощь" {
				vk.Method.SendMessage(helpmsg, v.Object.GetInt64("peer_id"))
			}
		}
	}
}
