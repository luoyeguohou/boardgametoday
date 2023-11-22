package main

import simplejson "github.com/bitly/go-simplejson"

type ModuleToufu struct {
	ply *Player
}

func (m *ModuleToufu) OnInit() {

}

func (m *ModuleToufu) OnExit() {

}

func (m *ModuleToufu) OnSec() {

}

func (m *ModuleToufu) ProcessMsg(cmd int, jParam *simplejson.Json) {
	if cmd == Toufu_Choose_Player {
		ServiceMGR.ServiceToufu.ChoosePlayer(ServiceMGR.ServiceRoom.FindRoom(m.ply).RoomID, jParam.Get("choose").MustInt())
	}
}
