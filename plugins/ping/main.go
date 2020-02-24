package main

import "../basic"

func Init(p *basic.Plugin) {
	p.Cmds = append(p.Cmds, "пинг", "жив?")
	p.Desc = "проверить работоспособность бота"
}

func Run(p *basic.Plugin) basic.APIError {
	p.VK.Method.SendMessage(":: Понг", p.Values.PeerID)

	return basic.APIError{}
}
