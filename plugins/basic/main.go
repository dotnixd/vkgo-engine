package basic

import (
	"database/sql"
	"log"

	"../../vkapi"
	"github.com/valyala/fastjson"
)

// Plugin описывает плагин
type Plugin struct {
	Cmds   []string
	Desc   string
	Data   *fastjson.Value
	VK     *vkapi.VK
	DB     *sql.DB
	Values struct {
		Args   []string
		PeerID int64
		FromID int64
	}
}

// APIError описывает ошибку API
type APIError struct {
	Service string
	Method  string
	Code    int
}

// ErrorNoArguments описывает ошибку недостатка аргументов
var ErrorNoArguments = APIError{
	Service: "bot",
	Method:  "no_enough_arguments",
	Code:    1,
}

// GenerateAPIError генерирует basic.APIError из JSON-представления
func GenerateAPIError(resp []byte) APIError {
	var p fastjson.Parser

	json, err1 := p.Parse(string(resp))
	Check(err1)

	if json.Exists("error") {
		err := json.Get("error")
		var apierror APIError

		for _, v := range err.GetArray("request_params") {
			if string(v.GetStringBytes("key")) == "method" {
				apierror.Method = string(v.GetStringBytes("value"))
			}
		}

		apierror.Service = "vk"
		apierror.Code = err.GetInt("error_code")

		return apierror
	}

	return APIError{}
}

// Check обработка исключений
func Check(err error, die ...bool) {
	if err != nil {
		if len(die) == 1 {
			log.Fatalln(err)
		} else {
			log.Println("Некритичная ошибка:", err)
		}
	}
}
