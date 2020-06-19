package main

import (
	"../../utils"
	"../basic"

	"github.com/valyala/fastjson"
)

func Init(p *basic.Plugin) {
	p.Cmds = append(p.Cmds, "онлайн", "оффлайн")
	p.Desc = "получить список участников которые онлайн/оффлайн"
}

func Run(p *basic.Plugin) basic.APIError {
	var parser fastjson.Parser

	members_nojson := p.VK.Method.GetConversationMembers(p.Values.PeerID)
	errapi := basic.GenerateAPIError(members_nojson)
	if errapi.Code != 0 {
		return errapi
	}

	members, err := parser.Parse(string(members_nojson))
	basic.Check(err)

	var online string
	var k string

	for _, v := range members.Get("response").GetArray("profiles") {
		if string([]rune(p.Values.Args[0])[1:]) == "онлайн" {
			k = "онлайн"
			if v.GetInt("online") == 1 {
				online += "\n" + "- [id" + utils.Int64ToString(v.GetInt64("id")) + "|" + string(v.GetStringBytes("first_name")) + " " + string(v.GetStringBytes("last_name")) + "]"
			}
		} else {
			if v.GetInt("online") == 0 {
				k = "оффлайн"
				online += "\n" + "- [id" + utils.Int64ToString(v.GetInt64("id")) + "|" + string(v.GetStringBytes("first_name")) + " " + string(v.GetStringBytes("last_name")) + "]"
			}
		}
	}

	p.VK.Method.SendMessage(":: Участники "+k+":"+online, p.Values.PeerID)

	return basic.APIError{}
}
