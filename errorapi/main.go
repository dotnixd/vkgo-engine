package errorapi

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"../plugins/basic"
	"../utils"
	"../vkapi"
	"github.com/valyala/fastjson"
)

// обработка ошибок
func _check(err error, die ...bool) {
	if err != nil {
		if len(die) == 1 {
			log.Fatalln(err)
		} else {
			log.Println("Некритичная ошибка:", err)
		}
	}
}

// ErrorHandler обработчик исключений API
type ErrorHandler struct {
	VK         *vkapi.VK
	Dictionary *fastjson.Value
}

// Init инициализирует структуру ErrorHandler
func (e *ErrorHandler) Init(vk *vkapi.VK, filename string) {
	log.Println(":: Чтение словаря ошибок")
	file, err := ioutil.ReadFile(filename)
	_check(err, true)

	var p fastjson.Parser
	data, err := p.Parse(string(file))
	_check(err)

	e.Dictionary = data
	e.VK = vk
}

// Handle обрабатывает исключения API
func (e *ErrorHandler) Handle(p *basic.Plugin, err *basic.APIError) {
	if err.Method == "" || err.Code == 0 {
		return
	}

	switch err.Service {
	case "", "vk":
		keys := strings.Split(err.Method, ".")
		section := e.Dictionary.Get("vk").Get(keys[0])
		if section.Exists(keys[1]) {
			method := section.GetArray(keys[1])
			for _, value := range method {
				if value.GetInt("code") == err.Code {
					e.VK.Method.SendMessage(":: Исключение: "+string(value.GetStringBytes("msg")),
						p.Values.PeerID)
					return
				}
			}
		}
	case "bot":
		section := e.Dictionary.Get("bot")
		if section.Exists(err.Method) {
			e.VK.Method.SendMessage(":: "+string(section.Get(err.Method).GetStringBytes("msg")), p.Values.PeerID)
			return
		}
	}

	msg := ":: Необработанное исключение\n"
	msg += "- Сервис: "
	if err.Service == "" {
		msg += "vk"
	} else {
		msg += err.Service
	}
	msg += "\n- Метод: " + err.Method + "\n"
	msg += "- Код ошибки: " + strconv.Itoa(err.Code) + "\n"
	msg += "- Команда: \"" + string(p.Data.GetStringBytes("text")) + "\"\n"
	msg += "- Peer_ID чата: " + utils.Int64ToString(p.Values.PeerID)

	var suffix string

	if e.Dictionary.Exists("logger_id") {
		e.VK.Method.SendMessage(msg, e.Dictionary.GetInt64("logger_id"))
		suffix = "\n\nОтчет отправлен разработчику"
	} else {
		suffix = "\n\nПожалуйста, перешлите этот лог разработчику"
	}

	e.VK.Method.SendMessage(msg+suffix, p.Values.PeerID)
}
