package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type ToufuPlayer struct {
	Gold      int
	Indentity int
	Player    *Player
}

func (p *ToufuPlayer) GetSImap() SImap {
	var ret = SImap{
		"Gold":     p.Gold,
		"Identity": p.Indentity,
		"IsPlayer": p.Player == nil,
	}

	if p.Player != nil {
		ret["Name"] = p.Player.NickName
	}

	return ret
}

type GameToufu struct {
	Room      *Room
	Players   [8]*ToufuPlayer
	PrinceIdx int
	PlayerNum int
}

type ServiceToufu struct {
	Rand  *rand.Rand
	Games map[int]*GameToufu
}

func (s *ServiceToufu) OnInit() {
	s.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	s.Games = make(map[int]*GameToufu)
}

func (s *ServiceToufu) OnPlyInit() {

}

func (s *ServiceToufu) OnSec() {
}

func (s *ServiceToufu) OnPlyExit(ply *Player) {
	// 玩家有参与某局游戏，销毁该局游戏
	room := ServiceMGR.ServiceRoom.FindRoom(ply)

	if room == nil {
		return
	}

	for _, game := range s.Games {
		if game == nil {
			continue
		}
		if game.Room == room {
			ServiceMGR.ServiceRoom.BroadcastRoom(room.RoomID, Toufu_Game_End, s.GetState(room.RoomID))
			s.Games[room.RoomID] = nil
		}
	}
}

func (s *ServiceToufu) GetRandomIdentity(princeIdx int) [8]int {
	tbl := []int{2, 3, 4, 5, 6, 7, 8}
	identities := [7]int{0, 0, 0, 0, 0, 0, 0}

	for i := 0; i < 7; i++ {
		randomNum := s.Rand.Intn(len(tbl))
		identities[i] = tbl[randomNum]
		tbl = append(tbl[:randomNum], tbl[randomNum+1:]...)
	}

	ret := [8]int{0, 0, 0, 0, 0, 0, 0, 0}
	for i := 0; i < 8; i++ {
		if i < princeIdx {
			ret[i] = identities[i]
		} else if i == princeIdx {
			ret[i] = 1
		} else if i > princeIdx {
			ret[i] = identities[i-1]
		}
	}

	return ret
}

func (s *ServiceToufu) GetState(roomID int) SImap {
	if s.Games[roomID] == nil {
		return nil
	}

	ret := make([]SImap, 0)
	for _, ply := range s.Games[roomID].Players {
		ret = append(ret, ply.GetSImap())
	}

	return SImap{"states": ret}
}

func (s *ServiceToufu) InitGame(room Room) int {
	// 游戏是否有效
	if s.Games[room.RoomID] != nil {
		return -1
	}
	// 随机第一局的
	s.Games[room.RoomID] = &GameToufu{Room: &room, PrinceIdx: 0, PlayerNum: 0}

	for i := 0; i < 8; i++ {
		if room.Players[i] != nil {
			s.Games[room.RoomID].PlayerNum++
		}
	}

	var randomIdentity = s.GetRandomIdentity(0)
	s.Games[room.RoomID].PrinceIdx++
	if s.Games[room.RoomID].PrinceIdx >= s.Games[room.RoomID].PlayerNum {
		s.Games[room.RoomID].PrinceIdx = 0
	}
	for i := 0; i < 8; i++ {
		s.Games[room.RoomID].Players[i] = &ToufuPlayer{
			Gold:      0,
			Player:    room.Players[i],
			Indentity: randomIdentity[i],
		}
	}
	state := s.GetState(room.RoomID)
	if state == nil {
		return -1
	}
	// 给房间里所有人播报游戏开始
	ServiceMGR.ServiceRoom.BroadcastRoom(room.RoomID, Start_Game, state)

	return 0
}

func (s *ServiceToufu) ChoosePlayer(roomID int, chooseIdx int) int {
	fmt.Println("choose player!!!")
	// 游戏是否有效
	if s.Games[roomID] == nil {
		fmt.Println("游戏无效")
		return -1
	}
	// 选人是否有效
	if chooseIdx < 0 || chooseIdx > 7 || s.Games[roomID].Players[chooseIdx].Indentity == 1 {
		fmt.Println("选人无效")
		return -1
	}
	// 结算金币

	goldGet := make([]int, 0)
	posGetGold := make([]int, 0)
	camp := ToufuCfgs[s.Games[roomID].Players[chooseIdx].Indentity].camp
	fmt.Println("翻牌的人的身份：" + strconv.Itoa(camp))

	for idx, ply := range s.Games[roomID].Players {
		fmt.Println("遍历的人的身份：" + strconv.Itoa(ToufuCfgs[ply.Indentity].winCamp))
		if ply.Player != nil && ToufuCfgs[ply.Indentity].winCamp == camp {
			ply.Gold += ToufuCfgs[ply.Indentity].winGold
			goldGet = append(goldGet, ToufuCfgs[ply.Indentity].winGold)
			posGetGold = append(posGetGold, idx)
		}
	}

	result := SImap{
		"chooseIdentity": s.Games[roomID].Players[chooseIdx].Indentity,
		"goldGet":        goldGet,
		"posGetGold":     posGetGold,
	}

	if s.Games[roomID].Players[chooseIdx].Player != nil {
		result["chooseName"] = s.Games[roomID].Players[chooseIdx].Player.NickName
	} else {
		result["chooseName"] = "空座位" + strconv.Itoa(chooseIdx+1)
	}

	// 随机下一局的
	var randomIdentity = s.GetRandomIdentity(s.Games[roomID].PrinceIdx)

	s.Games[roomID].PrinceIdx++
	if s.Games[roomID].PrinceIdx >= s.Games[roomID].PlayerNum {
		s.Games[roomID].PrinceIdx = 0
	}

	for i := 0; i < 8; i++ {
		s.Games[roomID].Players[i].Indentity = randomIdentity[i]
	}
	// 如果游戏结束，销毁游戏本身
	state := s.GetState(roomID)
	if state == nil {
		return -1
	}
	// 给所有人播报下一局的身份与这局的结果
	ServiceMGR.ServiceRoom.BroadcastRoom(roomID, Toufu_Turn_End, s.GetState(roomID).Merge(result))

	for _, ply := range s.Games[roomID].Players {
		if ply.Gold >= 7 {
			ServiceMGR.ServiceRoom.BroadcastRoom(roomID, Toufu_Game_End, state)
			s.Games[roomID] = nil
		}
	}
	return 0
}
