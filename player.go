package main

import (
	"encoding/json"
	"fmt"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"
)

type Player struct {
	NickName string
	HeadPic  string
	Conn     *websocket.Conn
	Modules  []Module
}

func (ply *Player) ProcessMsg(msg []byte) {
	jmsg, err := simplejson.NewJson(msg)
	if err != nil {
		return
	}
	v_msg, err := jmsg.Array()
	if err != nil || len(v_msg) == 0 || len(v_msg) > 2 {
		return
	}
	cmd := jmsg.GetIndex(0).MustInt()
	jParam := jmsg.GetIndex(1)

	for _, m := range ply.Modules {
		m.ProcessMsg(cmd, jParam)
	}

	fmt.Println(cmd)
	fmt.Println(jParam)
}

func (ply *Player) ToSImap() SImap {
	return SImap{
		"NickName": ply.NickName,
		"HeadPic":  ply.HeadPic,
	}
}

func (ply *Player) OnInit() {
	ply.Modules = make([]Module, 0)
	ply.Modules = append(ply.Modules, &ModuleRoom{ply: ply})
	ply.Modules = append(ply.Modules, &ModuleToufu{ply: ply})
	ply.Modules = append(ply.Modules, &ModuleHeartBeat{ply: ply})
	ply.NickName = "落叶过后"
	TimerRoutine(time.Second, func() { ply.onSec() }, nil)

}

func (ply *Player) OnExit() {

}

func (ply *Player) onSec() {
	for _, m := range ply.Modules {
		m.OnSec()
	}
}

func (ply *Player) SendMsg(cmd int, msg SImap) {
	ply.Conn.WriteMessage(1, Val2JsonByte([]interface{}{cmd, msg}))
}

type Module interface {
	OnInit()
	ProcessMsg(cmd int, jParam *simplejson.Json)
	OnExit()
	OnSec()
}

type SImap map[string]interface{}

func (si SImap) Merge(other SImap) SImap {
	for k, v := range other {
		si[k] = v
	}
	return si
}

func Val2JsonByte(val interface{}) []byte {
	b_json, err := json.Marshal(val) //Val2sjson(val).MarshalJSON()
	if err == nil {
		return b_json
	}
	return []byte(err.Error())
}
