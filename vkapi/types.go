package vkapi

import (
	"github.com/valyala/fastjson"
)

// LongPollEvent описывает структуру события
type LongPollEvent struct {
	TS      string
	Updates []LongPollEventUpdate
}

// LongPollEventUpdate событие, но поближе
type LongPollEventUpdate struct {
	Type   string
	Object *fastjson.Value
}
