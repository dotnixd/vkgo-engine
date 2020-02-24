package utils

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"unicode"
)

// ExtractCmd извлекает команду из сообщения
func ExtractCmd(cmd string) string {
	runes := []rune(cmd)
	ret := ""
	for k, v := range runes {
		if k == 0 {
			continue
		}
		if unicode.IsSpace(v) {
			return ret
		}
		ret += string(v)
	}

	return ret
}

// Int64ToString int64 в string
func Int64ToString(val int64) string {
	return strconv.FormatInt(val, 10)
}

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

// Post функция для POST запросов к серверам
func Post(urll string, v map[string]string) []byte {
	data := ""

	var is bool = true

	for i, v := range v {
		if is {
			data = data + i + "=" + url.QueryEscape(v)
			is = false
			continue
		}

		data = data + "&" + i + "=" + url.QueryEscape(v)
	}

	r := bytes.NewReader([]byte(data))

	resp, err := http.Post(urll, "application/x-www-form-urlencoded", r)
	_check(err)

	if err != nil {
		return []byte{}
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	_check(err)
	return body
}

// Get функция для GET запросов к серверам
func Get(urll string, v map[string]string) []byte {
	data := urll + "?"

	var is bool = true

	for i, v := range v {
		if is {
			data = data + i + "=" + url.QueryEscape(v)
			is = false
			continue
		}

		data = data + "&" + i + "=" + url.QueryEscape(v)
	}

	resp, err := http.Get(data)
	_check(err)

	if err != nil {
		return []byte{}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	_check(err)
	return body
}
