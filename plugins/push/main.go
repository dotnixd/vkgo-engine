package main

import (
	"strings"

	"../../utils"
	"../basic"
	"github.com/valyala/fastjson"
)

func Init(p *basic.Plugin) {
	p.Cmds = append(p.Cmds, "пуш", "объявление")
	p.Desc = "упомянуть всех участников беседы"
}

func Run(p *basic.Plugin) basic.APIError {
	members_nojson := p.VK.Method.GetConversationMembers(p.Values.PeerID)
	errapi := basic.GenerateAPIError(members_nojson)
	if errapi.Code != 0 {
		return errapi
	}

	var parser fastjson.Parser
	members_json, err := parser.Parse(string(members_nojson))
	basic.Check(err)

	user_ids := members_json.Get("response").GetArray("profiles")

	suffix := ""

	for _, v := range user_ids {
		suffix += "[id" + utils.Int64ToString(v.GetInt64("id")) + "|ᅠ]"
	}

	sender, err := parser.Parse(string(p.VK.Method.UserGet(p.Values.FromID)))
	basic.Check(err)

	p.VK.Method.SendMessage(":: [id"+utils.Int64ToString(p.Values.FromID)+"|"+
		string(sender.GetArray("response")[0].GetStringBytes("first_name"))+
		"] хочет сделать объявление: "+strings.Join(p.Values.Args[1:], " ")+suffix,
		p.Values.PeerID, true)

	return basic.APIError{}
}

func main() {}
