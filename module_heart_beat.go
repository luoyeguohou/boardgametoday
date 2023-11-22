package main

import simplejson "github.com/bitly/go-simplejson"

type ModuleHeartBeat struct {
	ply *Player
}

func (m *ModuleHeartBeat) OnInit() {

}

func (m *ModuleHeartBeat) OnExit() {

}

func (m *ModuleHeartBeat) OnSec() {

}

func (m *ModuleHeartBeat) ProcessMsg(cmd int, jParam *simplejson.Json) {
	if cmd == Heart_Beat {
		m.ply.SendMsg(Heart_Beat, SImap{})
	}
}
