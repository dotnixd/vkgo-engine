package main

import (
	"strings"

	"../../utils"
	"../basic"
	"github.com/valyala/fastjson"
)

func Init(p *basic.Plugin) {
	p.Cmds = append(p.Cmds, "переведи", "перевод", "переводчик")
	p.Desc = "перевести русский текст на английский язык, если язык сообщения другой, то сообщение будет переведено на русский язык"
	p.ConfigPos = "translate"
}

func Run(p *basic.Plugin) basic.APIError {
	if len(p.Values.Args) < 2 {
		return basic.ErrorNoArguments
	}

	txt := strings.Join(p.Values.Args[1:], " ")

	lang_nojson := string(utils.Post("https://translate.yandex.net/api/v1.5/tr.json/detect", map[string]string{
		"key":  "trnsl.1.1.20191031T124411Z.cd9c82c2cc463cb9.1b5043a3a12489d25f21ccca1441a6f263b6322f", // string(p.Config.GetStringBytes("token")),
		"text": txt,
	}))

	if lang_nojson == "" {
		return basic.APIError{Service: "yandex-translate", Method: "detect", Code: 1}
	}

	var parser fastjson.Parser

	lang, err := parser.Parse(lang_nojson)
	basic.Check(err)

	var l string

	if lang.Exists("lang") {
		if string(lang.GetStringBytes("lang")) == "ru" {
			l = "en"
		} else {
			l = "ru"
		}
	} else {
		return basic.APIError{Service: "yandex-translate", Method: "detect", Code: 2}
	}

	text_nojson := string(utils.Post("https://translate.yandex.net/api/v1.5/tr.json/translate", map[string]string{
		"key":  "trnsl.1.1.20191031T124411Z.cd9c82c2cc463cb9.1b5043a3a12489d25f21ccca1441a6f263b6322f", // string(p.Config.GetStringBytes("token")),
		"text": txt,
		"lang": l,
	}))

	if text_nojson == "" {
		return basic.APIError{Service: "yandex-translate", Method: "translate", Code: 1}
	}

	text, err := parser.Parse(text_nojson)
	basic.Check(err)

	if text.Exists("text") {
		for _, v := range text.GetArray("text") {
			p.VK.Method.SendMessage(":: Перевод: "+string(v.GetStringBytes()), p.Values.PeerID)
		}
	} else {
		return basic.APIError{Service: "yandex-translate", Method: "translate", Code: 2}
	}

	return basic.APIError{}
}
