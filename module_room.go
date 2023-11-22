package main

import (
	simplejson "github.com/bitly/go-simplejson"
)

type ModuleRoom struct {
	ply    *Player
	logCnt int
}

func (m *ModuleRoom) OnInit() {

}

func (m *ModuleRoom) OnExit() {

}

func (m *ModuleRoom) OnSec() {

}

func (m *ModuleRoom) ProcessMsg(cmd int, jParam *simplejson.Json) {
	if cmd == Create_Room {
		m.CreateRoom()
	} else if cmd == Join_Room {
		m.JoinRoom(jParam.Get("roomID").MustInt())
	} else if cmd == Change_Site {
		m.ChangeSite(jParam.Get("roomID").MustInt(), jParam.Get("aimSite").MustInt())
	} else if cmd == Exis_Room {
		m.ExitRoom(jParam.Get("roomID").MustInt())
	} else if cmd == Start_Game {
		m.StartGame()
	} else if cmd == Change_Name {
		m.ChangeName(jParam.Get("name").MustString())
	}
}

func (m *ModuleRoom) StartGame() {
	ServiceMGR.ServiceToufu.InitGame(*ServiceMGR.ServiceRoom.FindRoom(m.ply))
}

func (m *ModuleRoom) CreateRoom() {
	roomID := ServiceMGR.ServiceRoom.CreateRoom(m.ply)
	state := ServiceMGR.ServiceRoom.GetRoomState(roomID)
	if state == nil {
		m.ply.SendMsg(Create_Room, SImap{"ret": -1})
		return
	}
	m.ply.SendMsg(Create_Room, SImap{"ret": 0, "state": state, "roomID": roomID, "isHost": true, "selfIdx": ServiceMGR.ServiceRoom.FindIdx(m.ply)})
}

func (m *ModuleRoom) JoinRoom(roomID int) {
	ret := ServiceMGR.ServiceRoom.JoinRoom(roomID, m.ply)
	if ret != 0 {
		m.ply.SendMsg(Join_Room, SImap{"ret": -1})
		return
	}
	state := ServiceMGR.ServiceRoom.GetRoomState(roomID)
	if state == nil {
		m.ply.SendMsg(Join_Room, SImap{"ret": -2})
		return
	}
	m.ply.SendMsg(Join_Room, SImap{"ret": 0, "state": state, "roomID": roomID, "isHost": false, "selfIdx": ServiceMGR.ServiceRoom.FindIdx(m.ply)})
}

func (m *ModuleRoom) ChangeSite(roomID int, aimSite int) {
	ret := ServiceMGR.ServiceRoom.ChangePos(roomID, m.ply, aimSite)
	if ret != 0 {
		m.ply.SendMsg(Change_Site, SImap{"ret": -1})
		return
	}
}

func (m *ModuleRoom) ExitRoom(roomID int) {
	ret := ServiceMGR.ServiceRoom.ExitRoom(roomID, m.ply)
	if ret != 0 {
		m.ply.SendMsg(Exis_Room, SImap{"ret": -1})
		return
	}
	m.ply.SendMsg(Exis_Room, SImap{"ret": 0})
}

func (m *ModuleRoom) ChangeName(name string) {
	m.ply.NickName = name
	m.ply.SendMsg(Change_Name, SImap{"ret": 0, "newName": name})
}
