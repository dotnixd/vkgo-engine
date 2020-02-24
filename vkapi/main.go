package vkapi

import (
	"log"

	"../utils"
	"github.com/valyala/fastjson"
)

// VK структура для VKApi
type VK struct {
	Token   string
	GroupID string
	Method  Method
}

// VKApi URL API сервера ВКонтакте
const VKApi = "https://api.vk.com/method/"

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

// Request функция для совершения запросов к VKApi
func (v *VK) Request(method string, args map[string]string) []byte {
	args["access_token"] = v.Token
	args["v"] = "5.92"
	return utils.Post(VKApi+method, args)
}

// Init инициализирует структуру VK
func (v *VK) Init(token, groupid string) {
	v.Token = token
	v.GroupID = groupid
	v.Method = Method {VK: v}
}

// ListenLongPoll получает последние события
func (v *VK) ListenLongPoll() LongPollEvent {
	var p fastjson.Parser
	longpollevent, err := p.Parse(string(v.Request("groups.getLongPollServer", map[string]string{"group_id": v.GroupID})))
	_check(err)
	if err != nil {
		return LongPollEvent{}
	}

	l := longpollevent.Get("response")

	event := string(utils.Get(
		string(l.GetStringBytes("server")),
		map[string]string{
			"act":  "a_check",
			"wait": "25",
			"key":  string(l.GetStringBytes("key")),
			"ts":   string(l.GetStringBytes("ts"))}))

	return v.ParseLongPollEvent(string(event))
}

// ParseLongPollEvent переводит JSON-представление события в структуру LongPollEvent
func (v *VK) ParseLongPollEvent(event string) LongPollEvent {
	var e LongPollEvent
	var p fastjson.Parser
	val, err := p.Parse(event)
	_check(err, true)

	e.TS = string(val.GetStringBytes("ts"))
	upd := val.GetArray("updates")
	for _, k := range upd {
		var u LongPollEventUpdate
		u.Type = string(k.GetStringBytes("type"))
		u.Object = k.Get("object")
		e.Updates = append(e.Updates, u)
	}

	return e
}
