package main

import (
	"math/rand"
	"strings"
	"time"

	"../../utils"
	"../basic"
	"github.com/valyala/fastjson"
)

func Init(p *basic.Plugin) {
	p.Cmds = append(p.Cmds, "кто")
	p.Desc = "выбрать случайного участника беседы"
	rand.Seed(time.Now().UTC().UnixNano())
}

func Run(p *basic.Plugin) basic.APIError {
	if len(p.Values.Args) <= 1 {
		return basic.ErrorNoArguments
	}

	members_nojson := p.VK.Method.GetConversationMembers(p.Values.PeerID)
	errapi := basic.GenerateAPIError(members_nojson)
	if errapi.Code != 0 {
		return errapi
	}

	var parser fastjson.Parser
	members_json, err := parser.Parse(string(members_nojson))
	basic.Check(err)

	user_ids := members_json.Get("response").GetArray("items")
	user_id := user_ids[rand.Intn(len(user_ids))].GetInt64("member_id")

	user_nojson := p.VK.Method.UserGet(user_id)

	errapi = basic.GenerateAPIError(user_nojson)
	if errapi.Code != 0 {
		return errapi
	}

	// p.VK.Method.SendMessage(fmt.Sprint(user_nojson), p.Values.PeerID)

	user_json, err := parser.Parse(string(user_nojson))
	basic.Check(err)

	user := user_json.GetArray("response")[0]
	p.VK.Method.SendMessage(":: Кто "+strings.Join(p.Values.Args[1:], " ")+"? Возможно это [id"+
		utils.Int64ToString(user_id)+"|"+string(user.GetStringBytes("first_name"))+" "+
		string(user.GetStringBytes("last_name"))+"]", p.Values.PeerID)

	return basic.APIError{}
}

func main() {}
