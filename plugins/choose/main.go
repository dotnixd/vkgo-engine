package main

import (
	"math/rand"
	"time"

	"../basic"
)

func Init(p *basic.Plugin) {
	p.Cmds = append(p.Cmds, "выбери")
	p.Desc = "выбрать случайное значение"

	rand.Seed(time.Now().UTC().UnixNano())
}

func Run(p *basic.Plugin) basic.APIError {
	if len(p.Values.Args) <= 3 {
		return basic.ErrorNoArguments
	}

	var chooses []string
	var str string

	for _, v := range p.Values.Args[1:] {
		if v == "или" {
			chooses = append(chooses, str)
			str = ""
		} else {
			str += " " + v
		}
	}

	if str != "" {
		chooses = append(chooses, str)
	}

	ch := chooses[rand.Intn(len(chooses))]

	p.VK.Method.SendMessage(":: Выбрано:"+ch, p.Values.PeerID)

	return basic.APIError{}
}
